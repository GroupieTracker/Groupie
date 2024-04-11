package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var lastMusic string
var spotifyTracks []string
var Track string
var Score int

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

func init() {
	rand.Seed(time.Now().Unix())

	loadSpotifyTracks("static/spotify_tracks.json")
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./pages/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ChangeMusic(w http.ResponseWriter, r *http.Request) {
Start:
	randomSpotifyTrack := spotifyTracks[rand.Intn(len(spotifyTracks))]
	if Track == randomSpotifyTrack {
		goto Start
	}
	Track = randomSpotifyTrack

	http.Redirect(w, r, "/goBlindTest?music="+randomSpotifyTrack, http.StatusSeeOther)
}
func VerifTrack(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	answer := r.Form.Get("answer")
	if answer == GetTitleLocal(Track) {
		Score++
		http.Redirect(w, r, "/changeMusic", http.StatusSeeOther)
		return
	}

	fmt.Println("Bonne réponse:", GetTitleLocal(Track), "/ Réponse input:", answer)

	http.Redirect(w, r, "/goBlindTest?music="+Track, http.StatusSeeOther)
}

func GetTitle(w http.ResponseWriter, r *http.Request) {
	trackID := extractTrackID(Track)

	spotifyClient := initSpotifyClient()

	track, err := spotifyClient.GetTrack(spotify.ID(trackID))
	if err != nil {
		http.Error(w, "Error retrieving track: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Titre de la piste:", track.Name)

	http.Redirect(w, r, "/goBlindTest?music="+Track, http.StatusSeeOther)
}

func GetTitleLocal(Spottrack string) string {
	trackID := extractTrackID(Spottrack)

	spotifyClient := initSpotifyClient()

	track, _ := spotifyClient.GetTrack(spotify.ID(trackID))

	return track.Name

}

func extractTrackID(spotifyLink string) string {
	parts := strings.Split(spotifyLink, "/")

	trackID := parts[len(parts)-1]

	trackID = strings.Split(trackID, "?")[0]

	return trackID
}

func GoBlindTest(w http.ResponseWriter, r *http.Request) {
	music := r.URL.Query().Get("music")
	if music == "" {
		randomSpotifyTrack := spotifyTracks[rand.Intn(len(spotifyTracks))]
		for randomSpotifyTrack == lastMusic {
			randomSpotifyTrack = spotifyTracks[rand.Intn(len(spotifyTracks))]
		}
		lastMusic = randomSpotifyTrack
		music = randomSpotifyTrack
	}

	tmpl, err := template.ParseFiles("./pages/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Track = music

	data := struct {
		TrackSpot string
		Score     int
	}{
		TrackSpot: Track,
		Score:     Score,
	}

	tmpl.Execute(w, data)
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

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/goBlindTest", GoBlindTest)
	http.HandleFunc("/changeMusic", ChangeMusic)
	http.HandleFunc("/getTitle", GetTitle)
	http.HandleFunc("/verifTrack", VerifTrack)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
