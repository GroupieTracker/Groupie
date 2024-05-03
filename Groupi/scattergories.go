package Groupi

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

func WsScattergories(w http.ResponseWriter, r *http.Request, timeForRound int, round int) {
	var err error
	db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()
	roomID := r.URL.Query().Get("room")
	room, ok := rooms[roomID]
	if !ok {
		room = &Room{
			ID:          roomID,
			Connections: make(map[*websocket.Conn]bool),
		}
		rooms[roomID] = room
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket: l 141", err)
		return
	}
	fmt.Println("----------------------------GAME----------------------------")
	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()
	iDCreatorOfRoom, err := GetRoomCreatorID(db, roomID)
	if err != nil {
		log.Println("Error upgrading to WebSocket: l 151", err)
		return
	}
	roomIDInt, _ := strconv.Atoi(roomID)
	var answer []string
	var lettre string
	var tabAnswer [][]string
	var tabOpinion [][][]string
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		fmt.Println("Erreur lors de la récupération du cookie :", err)
		return
	}
	fmt.Println("Valeur du cookie:", cookie.Value)
	userID, _ := GetUserIDByUsername(db, cookie.Value)

	fmt.Println("----------------------------START----------------------------")
	for i := 0; i < round; i++ {
		fmt.Println("---------------------------------DEB DU TOUR---------------------------------")
		usersIDs, _ := GetUsersInRoom(db, roomID)
		userScores, err := GetUserScoresForRoom(db, usersIDs, roomIDInt)
		if err != nil {
			fmt.Println("Erreur lors de la get scores:", err)
			return
		}
		sort.Slice(userScores, func(i, j int) bool {
			return userScores[i][1] < userScores[j][1]
		})

		tabAnswer = [][]string{}
		tabOpinion = [][][]string{}

		//init round time+lettre
		stop := make(chan struct{})
		time.Sleep(1 * time.Second)
		if userID == iDCreatorOfRoom {
			sendScores(room, userScores)
			lettre = sendRandomLetter(room)
			go bouclTimer(room, timeForRound, stop)
		}
		//read message
		for {
			_, mess, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				mutex.Lock()
				delete(room.Connections, conn)
				mutex.Unlock()
				return
			}
			dataGame, err := parseEventData(mess)
			if err != nil {
				fmt.Println("Erreur lors de la conversion des données:", err)
			}
			if dataGame.Event == "end" {
				sendEndBut(room)

			} else if dataGame.Event == "catchBackData" {
				answer = dataGame.Data
				sendData(room, answer)
			} else if dataGame.Event == "allDataOnebyOne" {
				answer = dataGame.Data
				tabAnswer = append(tabAnswer, answer)
				if userID == iDCreatorOfRoom {
					if len(tabAnswer) == len(usersIDs) {
						sendScoreForResults(room, lettre, tabAnswer)
						for i := 0; i <= 5*len(tabAnswer); i++ {
							time.Sleep(1 * time.Second)
							sendTimer(room, (5*len(tabAnswer))-i)
						}
						sendEvent(room, "opinionBack")
					}
				}
			} else if dataGame.Event == "opinion" {
				fmt.Println("Opignion send")
				te := dataGame.Data2

				sendOpinion(room, te)
			} else if dataGame.Event == "allOpignionOnebyOne" {
				te := dataGame.Data2
				tabOpinion = append(tabOpinion, te)
				if userID == iDCreatorOfRoom {
					fmt.Println("data2", dataGame.Data2)
					if len(tabOpinion) == len(usersIDs) {
						fmt.Println("tabOpinion", tabOpinion)
						addScore(tabAnswer, lettre, roomIDInt, db, tabOpinion)
						break
					}
				}
			} else if dataGame.Event == "endTroun" {
				sendEvent(room, "fetchData")
			}
		}
		fmt.Println("---------------------------------DEB DU TOUR---------------------------------")
	}

}
