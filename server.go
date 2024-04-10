package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

var audioFiles []string

func init() {
	rand.Seed(time.Now().Unix())

	audioFiles, _ = findAudioFiles("static/audio")
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ChangeMusic(w http.ResponseWriter, r *http.Request) {
	randomAudio := audioFiles[rand.Intn(len(audioFiles))]

	http.Redirect(w, r, "/goBlindTest?music="+randomAudio, http.StatusSeeOther)
}

func GoBlindTest(w http.ResponseWriter, r *http.Request) {
	music := r.URL.Query().Get("music")
	if music == "" {
		music = audioFiles[rand.Intn(len(audioFiles))]
	}

	tmpl, err := template.ParseFiles("./pages/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, music)
}

func findAudioFiles(dir string) ([]string, error) {
	var audioFiles []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mp3" || filepath.Ext(file.Name()) == ".wav" {
			audioFiles = append(audioFiles, filepath.Join(dir, file.Name()))
		}
	}

	if len(audioFiles) == 0 {
		return nil, errors.New("no audio files found in the directory")
	}

	return audioFiles, nil
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
