package Groupi

import(
	"fmt"
	"log"
	"net/http"
	websocket"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

func boucl(room *Room) {
	for {
		time.Sleep(3 * time.Second)
		SendRandomLetter(room)
	}
}

func SendRandomLetter(room *Room) {
	letter := getRandomLetter()
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println(letter)
	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(letter))
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
}

func getRandomLetter() string {
	letters := [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "Y", "Z"}
	randomIndex := rand.Intn(len(letters))
	return letters[randomIndex]
}

func WsScattergories(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fmt")

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

	// Lancement de la goroutine pour l'envoi de lettres aléatoires
	go boucl(room)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			mutex.Lock()
			delete(room.Connections, conn)
			mutex.Unlock()
			return
		}
	}
}