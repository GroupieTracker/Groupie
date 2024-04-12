package Groupi

import(
	"fmt"
	"log"
	"net/http"
	websocket"github.com/gorilla/websocket"
	"math/rand"
	"time"
)
round := 5


//boucle de jeu {
	// envoit d'un lettre 
	// mettre a jours le crono 
					
	// fini ?? oui -->

// recuper les données de touts le monde 

// afiche les repose de tout le monde par categorie si elle ne sont pas nulles
//  ajjote le score en fonction dans la db

// next round 
// }
// if round ===0 ->fin
// affiche un tableau de score 
// buton replay --> to new lobby






func bouclTimer(room *Room) {
	time:= 10
	for {
		TimerForPablo(room , time)
		time=time-1
		time.Sleep(1 * time.Second)
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
func TimerForPablo(room *Room ,rftgyhu int) {

	mutex.Lock()
	defer mutex.Unlock()

	for conn := range room.Connections {
		tim := []byte(fmt.Sprintf("%d", rftgyhu))
		err := conn.WriteMessage(websocket.TextMessage, []byte(tim))
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

	
    go bouclTimer(room)

SendRandomLetter(room)
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            mutex.Lock()
            delete(room.Connections, conn)
            mutex.Unlock()
            return
        }
        
        // Afficher le message reçu
        fmt.Printf("Message reçu de la connexion %p : %s\n", conn, p)
if string(p)=="end" {
	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte("catchData"))
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
	
}
        // Si vous avez besoin du type du message, vous pouvez également l'afficher
        fmt.Printf("Type du message : %d\n", messageType)
    }
}
