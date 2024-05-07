package Groupi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	websocket "github.com/gorilla/websocket"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
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
var userName string
var trackTitle string
var playerInRoom []string
var inputAnswer string
var playersInRoomStruct []Player
var conWin bool = false
var myTime int
var timerRunning bool
var timerMutex sync.Mutex

func getRandomMusic() string {
	loadSpotifyTracks("static/assets/tracks/spotify_tracks.json")
Start:
	randomSpotifyTrack := spotifyTracks[rand.Intn(len(spotifyTracks))]
	if Track == randomSpotifyTrack {
		goto Start
	}
	return randomSpotifyTrack
}

type Player struct {
	Pseudo string
	Score  int
	Status bool
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
	trackTitle = GetTitle()

	for conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Error writing message:", err)
			delete(room.Connections, conn)
		}
	}
}

func addPlayer(username string) {
	for _, player := range playersInRoomStruct {
		if player.Pseudo == username {
			return
		}
	}
	newPlayer := Player{Pseudo: username, Score: 0, Status: true}
	playersInRoomStruct = append(playersInRoomStruct, newPlayer)
}

func initSpotifyClient() *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     "77ac44bd776c43f5b83101eb965ce2a0",
		ClientSecret: "d3caf0734bc84e439e02454b6510453e",
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

func bouclTimerBT(room *Room) {
	fmt.Println(len(room.Connections))
	timeForRound := 10
	for {
		myTime = timeForRound
		sendTimerBT(room, timeForRound)
		timeForRound = timeForRound - 1
		time.Sleep(1 * time.Second)
		if timeForRound < 0 {
			for i, _ := range playersInRoomStruct {
				playersInRoomStruct[i].Status = true
			}
			musicLock.Lock()
			SpotifyMusic(room)
			musicLock.Unlock()
			timeForRound = 10
		}
		TimerScore = timeForRound
	}
}

func sendTimerBT(room *Room, time int) {
	fmt.Println(playersInRoomStruct)
	var title string
	if Track != "" {
		title = trackTitle
	} else {
		title = "quoicoubebou des montagnes"
	}
	tabScore := struct {
		Event    string   `json:"event"`
		Time     int      `json:"time"`
		Title    string   `json:"title"`
		Username string   `json:"username"`
		Players  []Player `json:"players"`
		WinCond  bool     `json:"wincond"`
	}{
		Event:    "timer",
		Time:     time,
		Title:    title,
		Username: userName,
		Players:  playersInRoomStruct,
		WinCond:  conWin,
	}
	timerDataLock.Lock()
	defer timerDataLock.Unlock()

	for conn := range room.Connections {
		err := conn.WriteJSON(tabScore)
		if err != nil {
			log.Println("Error writing message:", err)
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

func addScoreStruct(username string, score int) {
	for i, player := range playersInRoomStruct {
		if player.Pseudo == username {
			playersInRoomStruct[i].Score += score
			playersInRoomStruct[i].Status = false
			return
		}
	}
}

func orderByScore(players []Player) {
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
}

func WsBlindTest(w http.ResponseWriter, r *http.Request, time int, nbRound int) {

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

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		fmt.Println("Erreur lors de la récupération du cookie :", err)
		return
	}

	userName := cookie.Value

	newPlayer := false
	for _, name := range playerInRoom {
		if userName == name {
			newPlayer = true
			break
		}
	}
	if !newPlayer {
		playerInRoom = append(playerInRoom, userName)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()

	timerMutex.Lock()
	if !timerRunning {
		go func() {
			timerRunning = true
			bouclTimerBT(room)
			timerRunning = false
		}()
	}
	timerMutex.Unlock()

	fmt.Println(playerInRoom)

	SpotifyMusic(room)

	for {
		_, mess, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			mutex.Lock()
			delete(room.Connections, conn)
			mutex.Unlock()
			return
		}
		dataGame, err := parseEventData(mess)
		if err != nil {
			log.Println("Error parsing message:", err)
			mutex.Lock()
			delete(room.Connections, conn)
			mutex.Unlock()
			return
		}

		fmt.Println(dataGame)

		addPlayer(dataGame.Username)

		if dataGame.Event == "answer" {
			fmt.Println("ouais ouais")
			fmt.Println("la réponse de:", dataGame.Username, " est:", dataGame.Answer)
			inputAnswer = dataGame.Answer
			if strings.ToLower(inputAnswer) == strings.ToLower(trackTitle) {
				for _, player := range playersInRoomStruct {
					fmt.Println(player.Status)
					if player.Status == true && player.Pseudo == dataGame.Username {
						addScoreStruct(dataGame.Username, myTime)
					}
					if player.Score == 100 {
						fmt.Println("le Gagnant est:", player.Pseudo)
						conWin = true
						return
					}
				}
				fmt.Println("c'est le bon titre")
			}
		}

		orderByScore(playersInRoomStruct)

	}

	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()
}
