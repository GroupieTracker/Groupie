package Groupi
import(
	"fmt"
	"log"
	"net/http"
    "sync"	
    websocket"github.com/gorilla/websocket"
)
type Room struct {
    ID          string
    Connections map[*websocket.Conn]bool
}

var (
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
    rooms       = make(map[string]*Room) // Carte pour stocker toutes les rooms
    mutex       = sync.Mutex{}           // Mutex pour la synchronisation lors de la gestion des connexions
)



func WsBlindTest(w http.ResponseWriter, r *http.Request) {
	// Récupère l'identifiant de la room à partir des paramètres de la requête
	roomID := r.URL.Query().Get("room")
	fmt.Println(roomID)
	if roomID == "" {
		roomID="blindTest"
		// http.Error(w, "Missing room parameter", http.StatusBadRequest)
		// return
	}
	fmt.Println(roomID)
	
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
			log.Println(err)
			return
		}
		defer conn.Close()
	
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
