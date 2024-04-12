package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login.html")
}

func Register(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/register.html")
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Extraire les données du formulaire
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("new_username")
	password := r.Form.Get("new_password")

	// Vérifier si l'utilisateur existe déjà
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM USER WHERE pseudo = ?", username)
	err = row.Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si l'utilisateur existe déjà, renvoyer une erreur
	if count > 0 {
		http.Error(w, "Nom d'utilisateur déjà utilisé", http.StatusBadRequest)
		return
	}

	// Ajouter l'utilisateur à la base de données
	_, err = db.Exec("INSERT INTO USER (pseudo, password) VALUES (?, ?)", username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page de connexion
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Extraire les données du formulaire
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Vérifier les informations de connexion dans la base de données
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var storedPassword string
	row := db.QueryRow("SELECT password FROM USER WHERE pseudo = ?", username)
	err = row.Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Vérifier si le mot de passe correspond
	if storedPassword != password {
		http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Rediriger l'utilisateur vers la page de succès
	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

func Success(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Connexion réussie!")
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/handle-register", HandleRegister)
	http.HandleFunc("/handle-login", HandleLogin)
	http.HandleFunc("/success", Success)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
