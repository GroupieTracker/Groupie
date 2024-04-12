package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	Groupi "Groupi/Groupi"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoBlindTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./pages/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tab := []string{"Bonjour", "mon", "ami"}
	// if str!="" {

	// 	tab = append(tab, str)
	// }
	tmpl.Execute(w, nil)
}

func GoGuessTheSong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./pages/guessTheSong.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r)
	})

	// URL pour le rendu de la page du blind test
	http.HandleFunc("/goBlindTest", func(w http.ResponseWriter, r *http.Request) {
		GoBlindTest(w, r)
	})

	http.HandleFunc("/goGuessTheSong", func(w http.ResponseWriter, r *http.Request) {
		GoGuessTheSong(w, r)
	})

	http.HandleFunc("/goBlindTest/webs", Groupi.WsBlindTest)
	http.HandleFunc("/goGuessTheSong/webs", Groupi.WsGuessTheSong)

	// Serveur de fichiers statiques
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
