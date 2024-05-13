package Groupi

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	websocket "github.com/gorilla/websocket"
)

func WsWaitingRoomGuessTheSong(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	defer conn.Close()
	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()
	iDCreatorOfRoom, err := GetRoomCreatorID(db, roomID)
	if err != nil {
		log.Println("Error upgrading to WebSocket : ", err)
		return
	}
	userNameOfCreator, _ := GetUsernameByID(db, iDCreatorOfRoom)
	roomIDInt, _ := strconv.Atoi(roomID)
	maxPlayer, err := GetMaxPlayersForRoom(db, roomIDInt)
	if err != nil {
		fmt.Println("Erreur lors de GetMaxPlayersForRoom:", err)
		return
	}
	for {
		usersIDs, _ := GetUsersInRoom(db, roomID)
		nbPlayer := len(usersIDs)
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
		fmt.Println(donnee)
		if err != nil {
			fmt.Println("Erreur lors de la conversion des données:", err)
			return
		}
		if donnee.Event == "newPlayer" {
			if nbPlayer >= maxPlayer {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
			db, err = sql.Open("sqlite3", "./Groupi/BDD.db")
			if err != nil {
				log.Fatal("Error opening database:", err)
			}
			defer db.Close()
			fmt.Println("newPLayer : ", donnee.Data[0])
			id, _ := GetUserIDByUsername(db, donnee.Data[0])
			err := AddRoomUser(db, roomIDInt, id)
			if err != nil {
				log.Println("Error reading message:", err)
			}
		} else if donnee.Event == "start" {
			sendStartSignal(room)
			break
		}
	}
}
