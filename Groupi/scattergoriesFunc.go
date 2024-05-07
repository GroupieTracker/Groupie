package Groupi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type BackData struct {
	Event string     `json:"event"`
	Data  []string   `json:"data"`
	Data2 [][]string `json:"data2"`
}

func parseEventData(data []byte) (*BackData, error) {
	var event1 BackData
	err := json.Unmarshal(data, &event1)
	if err == nil {
		return &event1, nil
	}

	return nil, fmt.Errorf("failed to parse data into BackData ")
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
			if timeactu <= -1 {
				sendEvent(room, "fetchData")
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func sendWaitingRoom(room *Room, nbPlayer int, maxPlayer int, username string) {
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
			log.Println("Error writing message SENDwAIITNGROOM:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
}
func sendOpinion(room *Room, dd [][]string) {
	tabscores := struct {
		Event string     `json:"event"`
		Opi   [][]string `json:"opi"`
	}{
		Event: "opinionForSend",
		Opi:   dd,
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
			log.Println("Error writing message SENDSCORE:", err)
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
			log.Println("Error writing message SEND RANDOM LETTRE:", err)
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
			log.Println("Error writing message SEND TIMER:", err)
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
			log.Println("Error writing message SENDSCORE:", err)
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
		Event   string `json:"event"`
		Nothing string `json:"nothing"`
	}{
		Event:   "start",
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
			log.Println("Error writing message SEND SATRT SIGNAL:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
}

func sendData(room *Room, data []string) {
	tabId := struct {
		Event string   `json:"event"`
		Data  []string `json:"data"`
	}{
		Event: "dataForSend",
		Data:  data,
	}
	tab, err := json.Marshal(tabId)
	if err != nil {
		fmt.Println("Erreur de marshalling JSON:", err)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, tab)
		if err != nil {
			log.Println("Error writing message SEND dATA:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
}

func sendEndBut(room *Room) {
	tabEnd := struct {
		Event string `json:"event"`
	}{
		Event: "endNow",
	}
	data, err := json.Marshal(tabEnd)
	if err != nil {
		fmt.Println("Erreur de marshalling JSON:", err)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error writing message SEND ENDBUT:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}

}
func sendEvent(room *Room, eve string) {
	if eve == "fetchData" {
		stop := make(chan struct{})
		go stopTimer(stop)
	}
	tabCatchData := struct {
		Event string `json:"event"`
		R     int    `json:"r"`
	}{
		Event: eve,
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
			log.Println("Error writing message SEND EVENT:", err)
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
			log.Println("Error writing message SENND SCOREFORRESULT:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}

}

func addScore(tabAnswer [][]string, lettre string, roomIDInt int, db *sql.DB, tabopinion [][][]string) {
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
					moy := 0.0
					for j := 0; j < len(tabopinion); j++ {
						sc, err := strconv.Atoi(tabopinion[j][i][y])
						if err != nil {
							fmt.Println("Erreur lors de la conversion des données:", err)
							return
						}
						moy += float64(sc)
					}
					moy = moy / float64(len(tabopinion))
					if moy > 0.5 {
						if unique {
							score += 2
						} else {
							score += 1
						}
					} else {
						score += 0
					}
				} else {
					score += 0
				}
			}
		}
		idactu, err := GetUserIDByUsername(db, tabAnswer[i][0])
		if err != nil {
			log.Println("Error GetUserIDByUsername: ", err)
			return
		}
		err = UpdateRoomUserScore(db, roomIDInt, idactu, score)
		if err != nil {
			fmt.Println("Erreur lors de la conversion des données:", err)
			return
		}
	}
}
