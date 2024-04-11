package Groupi

import(
	"fmt"
	"log"
	"net/http"
	websocket"github.com/gorilla/websocket"
	"math/rand"
	"time"
)


func boucl()  {

	time.Sleep(20 * time.Second)
		SendRandomLetter()

}


func SendRandomLetter() {
		letter := getRandomLetter()
		mutex.Lock()
		fmt.Println(	letter)
		for _, room := range rooms {
			for conn := range room.Connections {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(letter)); err != nil {
					log.Println("err : ",err)
					conn.Close()
					delete(room.Connections, conn)
				}
			}
		}
		mutex.Unlock()
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
		log.Println("dzudbzhda:",err)
		return
	}
	defer conn.Close()

	boucl()
	// Ajoute la connexion à la liste des connexions de la room
	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// En cas de déconnexion, supprime la connexion de la liste de la room
			mutex.Lock()
			delete(room.Connections, conn)
			mutex.Unlock()
			return
		}

		fmt.Println(string(p)) // Pour afficher le message reçu côté serveur

		// Diffuse le message à toutes les connexions de la room
		mutex.Lock()
		for conn := range room.Connections {
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				conn.Close()
				delete(room.Connections, conn)
			}
		}
		mutex.Unlock()
	}

}