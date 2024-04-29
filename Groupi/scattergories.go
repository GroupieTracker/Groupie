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

func addScore(tabAnswer [][]string, lettre string, roomIDInt int, db *sql.DB) {
	fmt.Println("tab", tabAnswer)
	nbCategories := 5
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

	if userID == iDCreatorOfRoom {

	// game
	if !isStarted {
		for {
			usersIDs, _ = GetUsersInRoom(db, roomID)
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
				db, err = sql.Open("sqlite3", "./Groupi/BDD.db")
				if err != nil {
				log.Fatal("Error opening database:", err)
					}
				defer db.Close()
				fmt.Println("newPLayer : " ,donnee.Data[0] )
				id , _:= GetUserIDByUsername(db, donnee.Data[0])
				err := AddRoomUser(db , roomIDInt , id)
				if err != nil {
					log.Println("Error reading message:", err)
				}
			}else if donnee.Event == "start" {
				isStarted = !isStarted
				sendStartSignal(room)
				break 
			}
		}
	} 
	if isStarted {
		fmt.Println("----------------------------START----------------------------")
		for i := 0; i < round; i++ {
			usersIDs,_=GetUsersInRoom(db, roomID)
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
					stopTimer(stop)
					endStart(room)

				} else if donnee.Event == "catchBackData" {

					answer = donnee.Data
					tabAnswer = append(tabAnswer, answer)
					fmt.Println(tabAnswer)
					if len(tabAnswer) == len(usersIDs) {
						if userID == iDCreatorOfRoom {
							addScore(tabAnswer, lettre, roomIDInt, db)
							break
						}
					}
				}
			}
		}
	}}
}
