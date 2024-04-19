package Groupi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
"sort"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type Room struct {
	ID          string
	Connections map[*websocket.Conn]bool
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rooms = make(map[string]*Room)
	mutex = sync.Mutex{}
)

type BackData struct {
	Event string   `json:"event"`
	Data  []string `json:"data"`
}

func parseEventData(data []byte) (*BackData, error) {
	var event BackData
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func getRandomLetter() string {
	letters := [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "Y", "Z"}
	randomIndex := rand.Intn(len(letters))
	return letters[randomIndex]
}
func bouclTimer(room *Room, timeForRound int) {
	timeactu := timeForRound
	for {
		sendTimer(room, timeactu)
		timeactu = timeactu - 1
		if timeactu <= 0 {
			endStart(room)
		}
		time.Sleep(1 * time.Second)
	}
}

func sendRandomLetter(room *Room) string {
	letter := getRandomLetter()
	tabLettre := struct {
		Event  string `json:"event"`
		Lettre string `json:"lettre"`
	}{
		Event:  "letter",
		Lettre: letter,
	}
	data, err := json.Marshal(tabLettre)
	if err != nil {
		fmt.Println("Erreur de marshalling JSON:", err)
		return "a"
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
	return letter
}
func sendTimer(room *Room, time int) {
	tabId := struct {
		Event string `json:"event"`
		Time  int    `json:"time"`
	}{
		Event: "timer",
		Time:  time,
	}
	data, err := json.Marshal(tabId)
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


func sendScores(room *Room, scores [][]string) {
	tabscores := struct {
		Event string `json:"event"`
		Scores  [][]string    `json:"scores"`
	}{
		Event: "scoresData",
		Scores:  scores,
	}
	data, err := json.Marshal(tabscores)
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
func sendId(room *Room, conn *websocket.Conn, userID int) {
	tabId := struct {
		Event string `json:"event"`
		Id    int    `json:"id"`
	}{
		Event: "id",
		Id:    userID,
	}
	data, err := json.Marshal(tabId)
	if err != nil {
		fmt.Println("Erreur de marshalling JSON:", err)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	if room.Connections[conn] {
		err := conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("Error writing message to connection:", err)
		}
	}

}



func WsScattergories(w http.ResponseWriter, r *http.Request, time int, round int, username string) {
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
	userID, err := GetUserIDByUsername(db, username)
	sendId(room, conn, userID)
	if err != nil {
		log.Println("Error upgrading to WebSocket: l156", err)
		return
	}
	roomIDInt, _ := strconv.Atoi(roomID)
	AddRoomUser(db, roomIDInt, userID)

	var answer []string
	var lettre string
	var tabAnswer [][]string
	var tabNul [][]string

	for i := 0; i < round; i++ {
		usersIDs, _ := GetUsersInRoom(db, roomID)
		userScores,err :=GetUserScoresForRoom(db , usersIDs , roomIDInt)
		if err != nil {
			fmt.Println("Erreur lors de la get scores:", err)
			return
		}
		sort.Slice(userScores, func(i, j int) bool {
			return userScores[i][1] < userScores[j][1]
		})
		fmt.Println(userScores)
		sendScores(room , userScores)

		tabAnswer=tabNul
		
		if err != nil {
			fmt.Println("Erreur lors de la conversion des données:", err)
			return
		}

		if userID == iDCreatorOfRoom {

			lettre = sendRandomLetter(room)
			go bouclTimer(room, time)
		}
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

			if donnee.Event == "end" {
				endStart(room)

			} else if donnee.Event == "catchBackData" {
				if userID == iDCreatorOfRoom {

					answer = donnee.Data
					tabAnswer = append(tabAnswer, answer)
					
					if len(tabAnswer) == len(usersIDs) {
						addScore(tabAnswer, lettre, roomIDInt, userID, db)
						break
					}
				}
			}
		}
	}
}

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
	fmt.Println("tab", tabAnswer,"|")
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
