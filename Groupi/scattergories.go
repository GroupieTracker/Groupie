package Groupi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	websocket"github.com/gorilla/websocket"
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

func sendRandomLetter(room *Room) {
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

func bouclTimer(room *Room) {
	timeForRound := 10
	// gere l'arre de la manche si le temsp arrive a 0
	for {
		sendTimer(room, timeForRound)
		timeForRound = timeForRound - 1
		time.Sleep(1 * time.Second)
	}
}

func sendTimer(room *Room, time int) {
	tabTime := struct {
		Event string `json:"event"`
		Time  int    `json:"time"`
	}{
		Event: "timer",
		Time:  time,
	}
	data, err := json.Marshal(tabTime)
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

func WsScattergories(w http.ResponseWriter, r *http.Request) {
  round := 5
	// Récupère l'identifiant de la room à partir des paramètres de la requête
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
    roomID = "petitBac"
	}
	// Vérifie si la room existe
	room, ok := rooms[roomID]
	if !ok {
    // Crée une nouvelle room si elle n'existe pas
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
  
  var answer []string
	for i := 0; i < round; i++ {
    
		//init start of round
		sendRandomLetter(room)
		//if id user === chef de la room pour evite les saut de timer {
		go bouclTimer(room)
		// Check if message
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
			} else if donnee.Event == "catchBackData" {
				answer  = donnee.Data
				fmt.Println(answer)
				

					// chef envoit l'id de de joueur qui doit envoyer c'est reponse
          

          } else if (1==1){ 
          // if idUser==idUser {
          // 	send(tabOfResult)
          // }

        //   var tabToSend []interface{}
		// tabToSend = append(tabToSend, "nomUser")
        //   tabToSend = append(tabToSend, answer)
        //   tabCatchData := struct {
        //     Event string `json:"event"`
        //     Answer   []interface    `json:"answer"`
        //   }{
        //     Event: "Answer",
        //     Answer:   tabToSend,
        //   }
        //   data, err := json.Marshal(tabCatchData)
        //   if err != nil {
        //     fmt.Println("Erreur de marshalling JSON:", err)
        //     return
        //   }
  
        //   for conn := range room.Connections {
        //     err := conn.WriteMessage(websocket.TextMessage, []byte(data))
        //     if err != nil {
        //       log.Println("Error writing message:", err)
        //       conn.Close()
        //       delete(room.Connections, conn)
        //     }
        //   }
					//envoyer ses données
				}
			}
		
	}
}

func displayAnswer()  {


  
}
