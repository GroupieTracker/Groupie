package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var lastMusic string
var spotifyTracks []string
var Lasttracks string

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
	if randomSpotifyTrack == Lasttracks {
		goto Start
	}
	Lasttracks = randomSpotifyTrack


	http.Redirect(w, r, "/goBlindTest?music="+randomSpotifyTrack, http.StatusSeeOther)
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
	tmpl.Execute(w, music)
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

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
