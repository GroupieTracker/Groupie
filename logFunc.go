package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	Groupi "Groupi/Groupi"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Redirect(w, r, "/lobby", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("static/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Redirect(w, r, "/lobby", http.StatusSeeOther)
		return
	}

	data := struct {
		Error string
	}{}
	tmpl, err := template.ParseFiles("static/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/lobby.html")
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)

	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	username := r.Form.Get("new_pseudo")
	email := r.Form.Get("new_email")
	password := r.Form.Get("new_password")

	var userError string

	if len(password) < 12 || !containsNumber(password) || !containsLetter(password) || !containsSpecialChar(password) {
		userError = "Le mot de passe doit contenir au moins 12 caractères, inclure au moins un chiffre, une lettre et un caractère spécial"
	} else {
		db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		defer db.Close()

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM USER WHERE pseudo = ?", username)
		err = row.Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}

		var count1 int
		row1 := db.QueryRow("SELECT COUNT(*) FROM USER WHERE email = ?", email)
		err = row1.Scan(&count1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}

		if count1 > 0 {
			userError = "Cette adresse mail est déjà utilisée"
		} else if count > 0 {
			userError = "Ce nom d'utilisateur est déjà utilisé"
		} else {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "MDP pas hashé", http.StatusInternalServerError)

			}

			_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			}
			fmt.Println(username)

			expiration := time.Now().Add(24 * time.Hour)
			cookieName := "auth_token"
			cookie := http.Cookie{
				Name:     cookieName,
				Value:    username,
				Expires:  expiration,
				HttpOnly: false,
			}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/lobby", http.StatusSeeOther)
		}
	}
	registerError(w, userError)

}

func registerError(w http.ResponseWriter, userError string) {
	tmpl, err := template.ParseFiles("static/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Error string
	}{
		Error: userError,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func containsNumber(s string) bool {
	r := regexp.MustCompile("[0-9]")
	return r.MatchString(s)
}

func containsLetter(s string) bool {
	r := regexp.MustCompile("[a-zA-Z]")
	return r.MatchString(s)
}

func containsSpecialChar(s string) bool {
	r := regexp.MustCompile(`[!@#$%^&*()_+{}[\]:;<>,.?/~]`)
	return r.MatchString(s)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)

	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	usernameOrEmail := r.Form.Get("username")
	password := r.Form.Get("password")

	db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	defer db.Close()

	var storedPassword string
	row := db.QueryRow("SELECT password FROM USER WHERE pseudo = ? OR email = ?", usernameOrEmail, usernameOrEmail)
	err = row.Scan(&storedPassword)
	if err != nil {
		loginError(w, "Nom d'utilisateur ou mot de passe incorrect")

	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		loginError(w, "Nom d'utilisateur ou mot de passe incorrect")

	}

	username, err := Groupi.GetUsernameByEmailOrUsername(db, usernameOrEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	fmt.Println(username)
	expiration := time.Now().Add(24 * time.Hour)
	cookieName := "auth_token"
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    username,
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/lobby", http.StatusSeeOther)

}

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return false
	}

	if cookie.Value != "" {
		return true
	}

	return false
}

func loginError(w http.ResponseWriter, userError string) {
	tmpl, err := template.ParseFiles("static/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Error string
	}{
		Error: userError,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:    "auth_token",
		Value:   "",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
