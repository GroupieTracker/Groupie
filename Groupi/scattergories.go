package Groupi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

func sendEvent(room *Room, envent string) {
	stop := make(chan struct{})
	go stopTimer(stop)
	tabCatchData := struct {
		Event string `json:"event"`
		R     int    `json:"r"`
	}{
		Event: envent,
		R:     -1,
	}
	data, err := json.Marshal(tabCatchData)
	if err != nil {
		fmt.Println("Erreur de marshalling JSON:", err)
		return
	}

	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}

}

func sendScoreForResults(room *Room, lettre string, tabAnswer [][]string) {
	tabScores := struct {
		Event  string     `json:"event"`
		Scores [][]string `json:"scores"`
	}{
		Event:  "resultsData",
		Scores: tabAnswer,
	}
	data, err := json.Marshal(tabScores)
	if err != nil {
		fmt.Println("Erreur de marshalling JSON:", err)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}

}

func addScore(tabAnswer [][]string, lettre string, roomIDInt int, db *sql.DB) {
	fmt.Println("tab", tabAnswer)
	nbCategories := len(tabAnswer[0]) - 1
	unique := true
	lettreMin := strings.ToLower(lettre)
	lettreMaj := strings.ToUpper(lettre)
	for i := 0; i < len(tabAnswer); i++ {
		score := 0
		for y := 1; y <= nbCategories; y++ {
			unique = true
			if tabAnswer[i][y] == "" {
				score += 0
			} else {
				mot := tabAnswer[i][y]
				if strings.HasPrefix(strings.ToLower(mot), lettreMin) || strings.HasPrefix(strings.ToUpper(mot), lettreMaj) {
					for o := 0; o < len(tabAnswer); o++ {
						if strings.ToLower(tabAnswer[i][y]) == strings.ToLower(tabAnswer[o][y]) && o != i {
							unique = false
						}
					}
					if unique {
						score += 2
					} else {
						score += 1
					}
				} else {
					score += 0
				}
			}
		}
		idactu, _ := GetUserIDByUsername(db, tabAnswer[i][0])
		err := UpdateRoomUserScore(db, roomIDInt, idactu, score)
		if err != nil {
			fmt.Println("Erreur lors de la conversion des données:", err)
			return
		}
	}
}

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
	defer conn.Close()
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
	var tabNul [][]string
	var tabOpinion [][][]string
	var tabAvis [][]string
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

		tabAnswer = tabNul

		//init round time+lettre
		stop := make(chan struct{})
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
			dataGame2, _ := parseEventData2(mess)
			if dataGame.Event == "end" {
				sendEvent(room, "fetchData")
			} else if dataGame.Event == "catchBackData" {
				answer = dataGame.Data
				sendData(room, answer)
			} else if dataGame.Event == "allDataOnebyOne" {
				answer = dataGame.Data
				tabAnswer = append(tabAnswer, answer)
				if userID == iDCreatorOfRoom {
					if len(tabAnswer) == len(usersIDs) {
						sendScoreForResults(room, lettre, tabAnswer)
						for i := 0; i < 5*len(tabAnswer); i++ {
							time.Sleep(1 * time.Second)
							sendTimer(room, (5*len(tabAnswer))-i)
						}
						sendEvent(room, "opinionBack")
					}
				}
			} else if dataGame.Event == "endRound" {
				tabAvis = append(tabAvis, dataGame.Data)
				if len(tabAvis) == len(usersIDs) {
					fmt.Println("tabAvis", tabAvis)
					// updatScore()
					break
				}
			} else if dataGame2.Event == "opinion" {
				tabOpinion = append(tabOpinion, dataGame2.Data)
				fmt.Print(dataGame2.Data)
				if len(tabOpinion) == len(usersIDs) {
					fmt.Println("tabOpinion", tabOpinion)
					// addScore(tabAnswer, lettre, roomIDInt, db)
				}
			}

		}

	}
}
