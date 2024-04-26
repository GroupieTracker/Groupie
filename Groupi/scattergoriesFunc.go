package Groupi

import(
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
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
func bouclTimer(room *Room, timeForRound int, stop <-chan struct{}) {
	timeactu := timeForRound
	for {
		select {
		case <-stop:
			return
		default:
			sendTimer(room, timeactu)
			timeactu = timeactu - 1
			if timeactu <= 0 {
				endStart(room)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func sendWaitingRoom(room *Room, nbPlayer int , maxPlayer int, username string) {
	var tab []any
	tab = append(tab, nbPlayer)
	tab = append(tab, maxPlayer)
	tab = append(tab, username)
	tabwaiting := struct {
		Event string `json:"event"`
		Data  []any  `json:"data"`
	}{
		Event: "waiting",
		Data:  tab,
	}
	data, err := json.Marshal(tabwaiting)
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
		Event  string     `json:"event"`
		Scores [][]string `json:"scores"`
	}{
		Event:  "scoresData",
		Scores: scores,
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

func stopTimer(stop chan<- struct{}) {
	stop <- struct{}{}
}


func sendStartSignal(room *Room) {
	tabscores := struct {
		Event  string     `json:"event"`
		Nothing string `json:"nothing"`
	}{
		Event:  "start",
		Nothing: "r",
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
