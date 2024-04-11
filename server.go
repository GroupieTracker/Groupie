package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "./login.html")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pseudo := r.Form.Get("pseudo")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", pseudo, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

func Success(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Connexion réussie!")
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/success", Success)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
