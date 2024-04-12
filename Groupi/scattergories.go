package Groupi

import(
	"fmt"
	"log"
	"net/http"
	websocket"github.com/gorilla/websocket"
	"math/rand"
	"time"
	"regexp"
)



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






func getRandomLetter() string {
	letters := [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "Y", "Z"}
	randomIndex := rand.Intn(len(letters))
	return letters[randomIndex]
}

func sendRandomLetter(room *Room) {
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





func bouclTimer(room *Room) {
	timeForRound:= 10
	// gere l'arre de la manche si le temsp arrive a 0
	for {
		sendTimer(room , timeForRound)
		timeForRound=timeForRound-1
		time.Sleep(1 * time.Second)
	}
}
func sendTimer(room *Room ,rftgyhu int) {

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
			}else{
				if 1==1 {
					
					re := regexp.MustCompile(`"([^"]+)"`)
					correspondances := re.FindAllStringSubmatch(string(p), -1)
					var tabOf []string
					for _, match := range correspondances {
						tabOf = append(tabOf, match[1])
					}
					
					
					userPseudo := "nomDuJoueur"
					tabOfResult := make([]string, len(tabOf)+1)
					copy(tabOfResult[1:], tabOf)
					tabOfResult[0] = userPseudo
					fmt.Println(tabOfResult)
					
					
					// chef envoit l'id de de joueur qui doit envoyer c'est reponse
				}else{
					//envoyer ses données 
					


				}




							
			
		}
	
}
}
// }
}
