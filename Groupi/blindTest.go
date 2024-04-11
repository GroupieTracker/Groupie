package Groupi
import(
	"fmt"
	"log"
	"net/http"
    "sync"	
    websocket"github.com/gorilla/websocket"
)
var (
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
    connections = make(map[*websocket.Conn]bool) // Map pour stocker toutes les connexions WebSocket
    mutex       = sync.Mutex{}                   // Mutex pour la synchronisation lors de la gestion des connexions
)

func WsBlindTest(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    // Ajoute la connexion à la liste des connexions
    mutex.Lock()
    connections[conn] = true
    mutex.Unlock()

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            // En cas de déconnexion, supprime la connexion de la liste
            mutex.Lock()
            delete(connections, conn)
            mutex.Unlock()
            return
        }

        fmt.Println(string(p)) // Pour afficher le message reçu côté serveur

        // Diffuse le message à toutes les connexions ouvertes
        mutex.Lock()
        for conn := range connections {
            if err := conn.WriteMessage(messageType, p); err != nil {
                log.Println(err)
                conn.Close()
                delete(connections, conn)
            }
        }
        mutex.Unlock()
    }
}