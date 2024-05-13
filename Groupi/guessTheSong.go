package Groupi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	websocket "github.com/gorilla/websocket"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var spotifyClientGTS *spotify.Client
var spotifyTracksGTS []string
var TrackGTS string = "https://open.spotify.com/embed/track/2uqYupMHANxnwgeiXTZXzd?&autoplay=1"
var AnswerGTS string
var TimerScoreGTS int
var firstUserGTS bool = true
var ActMusicGTS string = "https://open.spotify.com/embed/track/2uqYupMHANxnwgeiXTZXzd?&autoplay=1"
var musicLockGTS sync.Mutex
var timerDataLockGTS sync.Mutex
var userNameGTS string
var trackTitleGTS string
var playerInRoomGTS []string
var inputAnswerGTS string
var playersInRoomStructGTS []Player
var conWinGTS bool = false
var myTimeGTS int
var timerRunningGTS bool
var timerMutexGTS sync.Mutex
var timeFromdbGTS int
var songs = []Song{
	{"Michael Jackson", "Thriller"},
	{"Queen", "Bohemian Rhapsody"},
}

type LyricsResponse struct {
	Lyrics string `json:"lyrics"`
}

type Song struct {
	Artist string
	Title  string
}

type BackDataGTS struct {
	Event    string     `json:"event"`
	Data     []string   `json:"data"`
	Data2    [][]string `json:"data2"`
	Answer   string     `json:"answer"`
	Username string     `json:"username"`
}

func SpotifyMusicGTS(room *Room) {
	spotifyTracksGTS := getRandomMusic()
	TrackGTS = spotifyTracksGTS
	sendMusicGTS(room, getLyrics(), getLyrics())
}

func getLyrics() string {
	songs := []Song{
		{"Michael Jackson", "Thriller"},
		{"Queen", "Bohemian Rhapsody"},
		{"The Beatles", "Hey Jude"},
		{"Bob Dylan", "Like a Rolling Stone"},
		{"Elvis Presley", "Jailhouse Rock"},
		{"Stevie Wonder", "Superstition"},
		{"Led Zeppelin", "Stairway to Heaven"},
		{"David Bowie", "Space Oddity"},
		{"Nirvana", "Smells Like Teen Spirit"},
		{"Prince", "Purple Rain"},
		{"The Eagles", "Hotel California"},
		{"Pink Floyd", "Comfortably Numb"},
		{"Bruce Springsteen", "Born to Run"},
		{"Bob Marley", "No Woman, No Cry"},
		{"Madonna", "Like a Prayer"},
		{"Radiohead", "Creep"},
		{"Johnny Cash", "Ring of Fire"},
		{"Metallica", "Enter Sandman"},
		{"AC/DC", "Back in Black"},
		{"Guns N' Roses", "Sweet Child o' Mine"},
		{"The Beach Boys", "Good Vibrations"},
		{"The Doors", "Light My Fire"},
		{"Aretha Franklin", "Respect"},
		{"The Bee Gees", "Stayin' Alive"},
		{"The Clash", "London Calling"},
		{"Elton John", "Rocket Man"},
		{"Fleetwood Mac", "Go Your Own Way"},
		{"Aerosmith", "Dream On"},
		{"The Who", "My Generation"},
		{"Ray Charles", "What'd I Say"},
		{"The Ramones", "Blitzkrieg Bop"},
		{"Beyoncé", "Crazy in Love"},
		{"Frank Sinatra", "My Way"},
		{"Eminem", "Lose Yourself"},
		{"Coldplay", "Fix You"},
		{"Whitney Houston", "I Will Always Love You"},
		{"Oasis", "Wonderwall"},
		{"Green Day", "American Idiot"},
		{"The Police", "Every Breath You Take"},
		{"Red Hot Chili Peppers", "Under the Bridge"},
		{"Kanye West", "Gold Digger"},
		{"Bill Withers", "Ain't No Sunshine"},
		{"Jay-Z", "Empire State of Mind"},
		{"Jimi Hendrix", "Purple Haze"},
		{"OutKast", "Hey Ya!"},
		{"The White Stripes", "Seven Nation Army"},
		{"Black Sabbath", "Paranoid"},
		{"Santana", "Smooth"},
		{"The Temptations", "My Girl"},
		{"Simon & Garfunkel", "Bridge Over Troubled Water"},
		{"The Supremes", "Stop! In the Name of Love"},
		{"Cream", "Sunshine of Your Love"},
		{"The Jackson 5", "I Want You Back"},
		{"Tina Turner", "What's Love Got to Do with It"},
		{"Deep Purple", "Smoke on the Water"},
		{"Adele", "Rolling in the Deep"},
		{"Justin Timberlake", "Cry Me a River"},
		{"Katy Perry", "Firework"},
		{"The Notorious B.I.G.", "Juicy"},
		{"Neil Young", "Heart of Gold"},
		{"The Cure", "Just Like Heaven"},
		{"Otis Redding", "(Sittin' On) The Dock of the Bay"},
		{"The Velvet Underground", "Heroin"},
		{"Drake", "Hotline Bling"},
		{"The Kinks", "You Really Got Me"},
		{"Muddy Waters", "Rollin' Stone"},
		{"The Byrds", "Mr. Tambourine Man"},
		{"R.E.M.", "Losing My Religion"},
		{"Frank Ocean", "Thinkin Bout You"},
		{"Lou Reed", "Walk on the Wild Side"},
		{"The Shirelles", "Will You Love Me Tomorrow"},
		{"Talking Heads", "Once in a Lifetime"},
		{"Johnny Cash", "I Walk the Line"},
		{"The Four Seasons", "Can't Take My Eyes Off You"},
		{"Eric Clapton", "Tears in Heaven"},
		{"James Brown", "I Got You (I Feel Good)"},
		{"Van Morrison", "Brown Eyed Girl"},
		{"Elvis Costello", "Pump It Up"},
		{"The Velvet Underground & Nico", "I'm Waiting for the Man"},
		{"Public Enemy", "Fight the Power"},
		{"The Everly Brothers", "Wake Up Little Susie"},
		{"John Lennon", "Imagine"},
		{"Creedence Clearwater Revival", "Fortunate Son"},
		{"Run-D.M.C.", "Walk This Way"},
		{"Lynyrd Skynyrd", "Sweet Home Alabama"},
		{"Beastie Boys", "Sabotage"},
		{"The Ronettes", "Be My Baby"},
		{"Rihanna", "Umbrella"},
		{"Sam Cooke", "A Change Is Gonna Come"},
		{"The Smiths", "There Is a Light That Never Goes Out"},
		{"Patti Smith", "Because the Night"},
		{"Pearl Jam", "Alive"},
		{"The Band", "The Weight"},
		{"Bon Jovi", "Livin' on a Prayer"},
		{"LCD Soundsystem", "All My Friends"},
	}

	rand.Seed(time.Now().UnixNano())

	index := rand.Intn(len(songs))
	selectedSong := songs[index]

	trackTitleGTS = selectedSong.Title

	fmt.Println("Artiste: ", selectedSong.Artist, " Nom de la musique: ", selectedSong.Title)

	url := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", selectedSong.Artist, selectedSong.Title)
	response, err := http.Get(url)
	if err != nil {
		return "feur"
	}
	defer response.Body.Close()

	var lyricsResponse LyricsResponse
	if err := json.NewDecoder(response.Body).Decode(&lyricsResponse); err != nil {
		return "feur"
	}

	return lyricsResponse.Lyrics
}

func sendMusicGTS(room *Room, musicURL string, lyrics string) {
	ActMusicGTS = musicURL
	lines := strings.Split(lyrics, "\n")
	if len(lines) > 1 {
		lyrics = strings.Join(lines[1:], "\n")
	}

	lyrics = strings.ReplaceAll(lyrics, "\n", "<br>")
	musicData := struct {
		Event  string `json:"event"`
		Music  string `json:"music"`
		Lyrics string `json:"lyrics"`
	}{
		Event:  "music",
		Music:  musicURL,
		Lyrics: lyrics,
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
			delete(room.Connections, conn)
		}
	}
}

func addPlayerGTS(username string) {
	for _, player := range playersInRoomStructGTS {
		if player.Pseudo == username {
			return
		}
	}
	newPlayer := Player{Pseudo: username, Score: 0, Status: true}
	playersInRoomStructGTS = append(playersInRoomStructGTS, newPlayer)
}

func initSpotifyClientGTS() *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     "77ac44bd776c43f5b83101eb965ce2a0",
		ClientSecret: "d3caf0734bc84e439e02454b6510453e",
		TokenURL:     spotify.TokenURL,
	}

	client := config.Client(context.Background())

	spotifyClientGTS := spotify.NewClient(client)

	return &spotifyClientGTS
}

func GetTitleGTS() string {
	var TitleTrack string
	trackID := extractTrackIDGTS(ActMusicGTS)

	spotifyClientGTS := initSpotifyClientGTS()

	track, err := spotifyClientGTS.GetTrack(spotify.ID(trackID))
	if err != nil {
		fmt.Println("Erreur lors de la récupération de la piste:", err)
		return ""
	}
	TitleTrack = track.Name

	return TitleTrack
}

func extractTrackIDGTS(spotifyLink string) string {
	parts := strings.Split(spotifyLink, "/")

	trackID := parts[len(parts)-1]

	trackID = strings.Split(trackID, "?")[0]

	return trackID
}

func bouclTimerGTS(room *Room, nbRoundDB int) {
	var nbRoundloop int
	fmt.Println(len(room.Connections))
	timeForRound := timeFromdbGTS
	fmt.Println(nbRoundDB)
	for {
		fmt.Println(nbRoundloop)
		myTimeGTS = timeForRound
		sendTimerBTGTS(room, timeForRound)
		timeForRound = timeForRound - 1
		time.Sleep(1 * time.Second)
		if timeForRound < 0 {
			for i, _ := range playersInRoomStructGTS {
				playersInRoomStructGTS[i].Status = true
			}
			musicLockGTS.Lock()
			SpotifyMusicGTS(room)
			musicLockGTS.Unlock()
			timeForRound = timeFromdbGTS
			nbRoundloop++
			if nbRoundloop >= nbRoundDB {
				conWinGTS = true
			}
		}
		TimerScoreGTS = timeForRound
	}

}

func sendTimerBTGTS(room *Room, time int) {
	//fmt.Println(playersInRoomStructGTS, " conWinGTS est : ", conWinGTS)
	var title string
	if TrackGTS != "" {
		title = trackTitleGTS
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
		Username: userNameGTS,
		Players:  playersInRoomStructGTS,
		WinCond:  conWinGTS,
	}
	timerDataLockGTS.Lock()
	defer timerDataLockGTS.Unlock()

	for conn := range room.Connections {
		err := conn.WriteJSON(tabScore)
		if err != nil {
			log.Println("Error writing message:", err)
			delete(room.Connections, conn)
		}
	}
}

func addScoreStructGTS(username string, score int) {
	for i, player := range playersInRoomStructGTS {
		if player.Pseudo == username {
			playersInRoomStructGTS[i].Score += score
			playersInRoomStructGTS[i].Status = false
			return
		}
	}
}

func orderByScoreGTS(players []Player) {
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
}

func WsGuessTheSong(w http.ResponseWriter, r *http.Request, time int, nbRound int, dif string) {

	fmt.Println("le temps est:", time, " le nbRound est: ", nbRound)
	timeFromdbGTS = time
	roomID := r.URL.Query().Get("room")
	fmt.Println(roomID)
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		log.Println("Error converting room ID to int:", err)
		return
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

	userNameGTS := cookie.Value

	newPlayer := false
	for _, name := range playerInRoomGTS {
		if userNameGTS == name {
			newPlayer = true
			break
		}
	}
	if !newPlayer {
		playerInRoomGTS = append(playerInRoomGTS, userNameGTS)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()

	timerMutexGTS.Lock()
	if !timerRunningGTS {
		go func() {
			timerRunningGTS = true
			bouclTimerGTS(room, nbRound)
			timerRunningGTS = false
		}()
	}
	timerMutexGTS.Unlock()

	fmt.Println(playerInRoomGTS)

	SpotifyMusicGTS(room)

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

		addPlayerGTS(dataGame.Username)
		db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
		if err != nil {
			log.Fatal("Error opening database:", err)
		}
		defer db.Close()
		if dataGame.Event == "answer" {
			fmt.Println("ouais ouais")
			fmt.Println("la réponse de:", dataGame.Username, " est:", dataGame.Answer)
			inputAnswerGTS = dataGame.Answer
			if strings.ToLower(inputAnswerGTS) == strings.ToLower(trackTitleGTS) {
				for _, player := range playersInRoomStructGTS {
					fmt.Println(player.Status)
					if player.Status == true && player.Pseudo == dataGame.Username {
						addScoreStructGTS(dataGame.Username, myTimeGTS)
						updatePlayerScores(db, playersInRoomStructGTS, roomIDInt)
					}
					if player.Score == 10 {
						fmt.Println("le Gagnant est:", player.Pseudo)
						conWinGTS = true
						return
					}
				}
				fmt.Println("c'est le bon titre")
			}
		}

		orderByScoreGTS(playersInRoomStructGTS)

	}

	mutex.Lock()
	room.Connections[conn] = true
	mutex.Unlock()
}
