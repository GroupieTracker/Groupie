package Groupi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"

	websocket "github.com/gorilla/websocket"
)

var spotifyClient *spotify.Client
var spotifyTracks []string
var Track string = "https://open.spotify.com/embed/track/2uqYupMHANxnwgeiXTZXzd?&autoplay=1"
var Answer string
var Score int
var TimerScore int
var firstUser bool = true
var ActMusic string = "https://open.spotify.com/embed/track/2uqYupMHANxnwgeiXTZXzd?&autoplay=1"
var musicLock sync.Mutex
var timerDataLock sync.Mutex

type Room struct {
	ID          string
	Connections map[*websocket.Conn]bool
	UsersCount  int
	TimerActive bool
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rooms = make(map[string]*Room) // Carte pour stocker toutes les rooms
	mutex = sync.Mutex{}           // Mutex pour la synchronisation lors de la gestion des connexions
)

func getRandomMusic() string {
	loadSpotifyTracks("Groupi/assets/spotify_tracks.json")
Start:
	randomSpotifyTrack := spotifyTracks[rand.Intn(len(spotifyTracks))]
	if Track == randomSpotifyTrack {
		goto Start
	}
	return randomSpotifyTrack
}

func SpotifyMusic(room *Room) {
	spotifyTracks := getRandomMusic()
	Track = spotifyTracks
	sendMusic(room, spotifyTracks)
}

func sendMusic(room *Room, musicURL string) {
	fmt.Print("ça envoie le paquet !")
	ActMusic = musicURL
	fmt.Println("musique envoyer:", musicURL)
	musicData := struct {
		Event string `json:"event"`
		Music string `json:"music"`
	}{
		Event: "music",
		Music: musicURL,
	}

	jsonData, err := json.Marshal(musicData)
	if err != nil {
		log.Println("Erreur de marshalling JSON:", err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
}

func initSpotifyClient() *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     "5406ffed7c33489b915321267e3ca75f",
		ClientSecret: "bbaefa27380049a39c30e53c9bddc0c6",
		TokenURL:     spotify.TokenURL,
	}

	client := config.Client(context.Background())

	spotifyClient := spotify.NewClient(client)

	return &spotifyClient
}

func GetTitle() string {
	var TitleTrack string
	trackID := extractTrackID(ActMusic)

	spotifyClient := initSpotifyClient()

	track, err := spotifyClient.GetTrack(spotify.ID(trackID))
	if err != nil {
		fmt.Println("Erreur lors de la récupération de la piste:", err)
		return ""
	}
	TitleTrack = track.Name

	return TitleTrack
}

func extractTrackID(spotifyLink string) string {
	parts := strings.Split(spotifyLink, "/")

	trackID := parts[len(parts)-1]

	trackID = strings.Split(trackID, "?")[0]

	return trackID
}

func bouclTimer(room *Room) {
	if len(room.Connections) == 1 {
		timeForRound := 10
		for {
			sendTimer(room, timeForRound)
			timeForRound = timeForRound - 1
			time.Sleep(1 * time.Second)
			if timeForRound < 0 {
				musicLock.Lock()
				SpotifyMusic(room)
				musicLock.Unlock()
				timeForRound = 10
			}
			TimerScore = timeForRound
		}
	}
}

func sendTimer(room *Room, time int) {
	var title string
	if Track != "" {
		title = GetTitle()
	} else {
		title = "Austin"
	}
	tabTime := struct {
		Event string `json:"event"`
		Time  int    `json:"time"`
		Title string `json:"title"`
	}{
		Event: "timer",
		Time:  time,
		Title: title,
	}
	data, err := json.Marshal(tabTime)
	if err != nil {
		fmt.Println("Erreur de marshalling JSON:", err)
		return
	}

	timerDataLock.Lock()
	defer timerDataLock.Unlock()

	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
}

func loadSpotifyTracks(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &spotifyTracks)
	if err != nil {
		log.Fatal(err)
	}
}

func WsBlindTest(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		roomID = "blindTest"
	}

	room, ok := rooms[roomID]
	if !ok {
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

	if firstUser == true {
		go bouclTimer(room)
	}

	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()
}
