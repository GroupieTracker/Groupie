package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	Groupi "Groupi/Groupi"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Redirect(w, r, "/lobby", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "./static/index.html")
}

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
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("new_pseudo")
	email := r.Form.Get("new_email")
	password := r.Form.Get("new_password")

	var userError string

	if len(password) < 12 || !containsNumber(password) || !containsLetter(password) || !containsSpecialChar(password) {
		userError = "Le mot de passe doit contenir au moins 12 caractères, inclure au moins un chiffre, une lettre et un caractère spécial"
	} else {
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

		var count1 int
		row1 := db.QueryRow("SELECT COUNT(*) FROM USER WHERE email = ?", email)
		err = row1.Scan(&count1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count1 > 0 {
			userError = "Cette adresse mail est déjà utilisé"
		} else if count > 0 {
			userError = "Ce nom d'utilisateur est déjà utilisé"
		} else {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "MDP pas hasher", http.StatusInternalServerError)
				return
			}

			_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/lobby", http.StatusSeeOther)
			return
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
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usernameOrEmail := r.Form.Get("username")
	password := r.Form.Get("password")

	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var storedPassword string
	row := db.QueryRow("SELECT password FROM USER WHERE pseudo = ? OR email = ?", usernameOrEmail, usernameOrEmail)
	err = row.Scan(&storedPassword)
	if err != nil {
		loginError(w, "Nom d'utilisateur ou mot de passe incorrect")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		loginError(w, "Nom d'utilisateur ou mot de passe incorrect")
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookieName := "auth_token"
	cookieValue := usernameOrEmail
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Expires:  expiration,
		HttpOnly: true,
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
	// Supprimer le cookie d'authentification en fixant sa date d'expiration à une date antérieure
	expiration := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:    "auth_token",
		Value:   "",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)

	// Rediriger l'utilisateur vers la page de connexion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GoBlindTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GOLobbyOfScattergories() {
	//Creation d'une nouvelle party
	//type petiti bac
	//REcupere L'id de l'useur et le mais en t'en que createur
}

func GoGuessTheSong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/guessTheSong.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func GoScattergories(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/scattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/lobby", Lobby)
	http.HandleFunc("/handle-register", HandleRegister)
	http.HandleFunc("/handle-login", HandleLogin)
	http.HandleFunc("/BlindTest/webs", Groupi.WsBlindTest)
	http.HandleFunc("/GuessTheSong/webs", Groupi.WsGuessTheSong)
	http.HandleFunc("/Scattergories/webs", Groupi.WsScattergories)
	http.HandleFunc("/logout", Logout)

	http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
		GoBlindTest(w, r)
	})
	http.HandleFunc("/GuessTheSong", func(w http.ResponseWriter, r *http.Request) {
		GoGuessTheSong(w, r)
	})
	http.HandleFunc("/Scattergories", func(w http.ResponseWriter, r *http.Request) {
		GoScattergories(w, r)
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("static/", fs))

	fsCSS := http.FileServer(http.Dir("static/css"))
	http.Handle("/static/css/", http.StripPrefix("/static/css/", fsCSS))

	fsPicture := http.FileServer(http.Dir("static/assets/pictures"))
	http.Handle("/static/assets/pictures/", http.StripPrefix("/static/assets/pictures/", fsPicture))

	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
