package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Map pour conserver la liste des clients connectés
var broadcast = make(chan Message)           // Channel pour diffuser les messages à tous les clients

// Configuration de la mise en forme des messages
type Message struct {
	Content string `json:"content"`
}

// Configuration de la mise en forme des upgradées
var upgrader = websocket.Upgrader{}

func main() {
	// Gestion des routes
	http.HandleFunc("/ws", handleConnections)

	// Démarrage de la goroutine pour diffuser les messages aux clients
	go handleMessages()

	// Démarrage du serveur
	log.Println("Serveur WebSocket démarré sur :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Erreur de démarrage du serveur: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Mise à niveau de la connexion HTTP à une connexion WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Fermer la connexion lorsque la fonction retourne
	defer ws.Close()

	// Ajouter le client à la liste des clients connectés
	clients[ws] = true

	for {
		var msg Message
		// Lire le message du client
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erreur de lecture du message: %v", err)
			delete(clients, ws) // Supprimer le client de la liste en cas d'erreur
			break
		}
		// Envoyer le message reçu à la goroutine de diffusion
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Récupérer le prochain message de la chaîne de diffusion
		msg := <-broadcast
		// Générer le code HTML pour afficher le message
		messageHTML := "<div>" + msg.Content + "</div>"
		// Ajouter le message à une variable globale contenant tous les messages
		allMessages += messageHTML
		// Actualiser la page avec tous les messages
		http.HandleFunc("/", handleRoot)
		http.ListenAndServe(":8000", nil)
	}
}