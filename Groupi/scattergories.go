package Groupi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)


func endStart(room *Room) {
	tabCatchData := struct {
		Event string `json:"event"`
		r     int    `json:"r"`
	}{
		Event: "fetchData",
		r:     -1,
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

func addScore(tabAnswer [][]string, lettre string, roomIDInt int, userID int, db *sql.DB) {
	fmt.Println("tab", tabAnswer, "|")
	// [[82=iduser O o o fez fezf]]

	unique := true
	for i := 0; i < len(tabAnswer); i++ {
		score := 0
		for y := 1; y <= 5; y++ {
			if string(tabAnswer[i][y]) == "" {
				score += 0
			} else {
				for o := 0; o < len(tabAnswer); o++ {
					if string(tabAnswer[i][y]) == string(tabAnswer[o][y]) && o != i {
						unique = false
					}
				}
				if unique {
					score += 2
				} else {
					score += 1
				}

			}
		}

		err := UpdateRoomUserScore(db, roomIDInt, userID, score)
		if err != nil {
			fmt.Println("Erreur lors de la conversion des données:", err)
			return
		}
	}
}

func WsScattergories(w http.ResponseWriter, r *http.Request, time int, round int , userID int) {
	isStarted := false
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
	defer conn.Close()
	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()	

	iDCreatorOfRoom, err := GetRoomCreatorID(db, roomID)
	if err != nil {
		log.Println("Error upgrading to WebSocket: l 151", err)
		return
		}	
		userNameOfCreator  , _:= GetUsernameByID(db ,iDCreatorOfRoom )
	roomIDInt, _ := strconv.Atoi(roomID)
	usersIDs, _ := GetUsersInRoom(db, roomID)

	var answer []string
	var lettre string
	var tabAnswer [][]string
	var tabNul [][]string
	maxPlayer, err := GetMaxPlayersForRoom(db, roomIDInt)
	if err != nil {
		fmt.Println("Erreur lors de GetMaxPlayersForRoom:", err)
		return
	}

	// game
	if !isStarted {
		for {
			nbPlayer := len(usersIDs)
			fmt.Println("usersIDs : ", usersIDs, nbPlayer)
			fmt.Println("userID : " , userID , ",creatorID : ",userNameOfCreator)
			sendWaitingRoom(room, nbPlayer, maxPlayer, userNameOfCreator)
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				mutex.Lock()
				delete(room.Connections, conn)
				mutex.Unlock()
				return
			}

			donnee, err := parseEventData(p)
			if err != nil {
				fmt.Println("Erreur lors de la conversion des données:", err)
				return
			}
			if donnee.Event=="newPlayer" {
				fmt.Println("newPLayer : " ,donnee.Data[0] )
				
			}else if donnee.Event == "start" {
				fmt.Printf("start")
				isStarted = !isStarted
				sendStartSignal(room)
				break 
			}
		}
	} 
	if isStarted {
		usersIDs,_=GetUsersInRoom(db, roomID)
		for i := 0; i < round; i++ {
			//Score
			fmt.Println("usersIDs : ", usersIDs)
			userScores, err := GetUserScoresForRoom(db, usersIDs, roomIDInt)

			if err != nil {
				fmt.Println("Erreur lors de la get scores:", err)
				return
			}
			sort.Slice(userScores, func(i, j int) bool {
				return userScores[i][1] < userScores[j][1]
			})
			sendScores(room, userScores)
			tabAnswer = tabNul
			if err != nil {
				fmt.Println("Erreur lors de la conversion des données:", err)
				return
			}
			//init round time+lettre
			stop := make(chan struct{})
			if userID == iDCreatorOfRoom {
				lettre = sendRandomLetter(room)
				go bouclTimer(room, time, stop)
			}
			//read message
			for {
				_, p, err := conn.ReadMessage()
				if err != nil {
					log.Println("Error reading message:", err)
					mutex.Lock()
					delete(room.Connections, conn)
					mutex.Unlock()
					return
				}

				donnee, err := parseEventData(p)
				if err != nil {
					fmt.Println("Erreur lors de la conversion des données:", err)
					return
				}
				fmt.Println("donne.data : ", donnee.Data)

				if donnee.Event == "end" {
					endStart(room)
					stopTimer(stop)

				} else if donnee.Event == "catchBackData" {
					answer = donnee.Data
					tabAnswer = append(tabAnswer, answer)
					fmt.Println(tabAnswer)
					if userID == iDCreatorOfRoom {
						if len(tabAnswer) == len(usersIDs) {
							addScore(tabAnswer, lettre, roomIDInt, userID, db)
							break
						}
					}
				}
			}
		}
	}
}
