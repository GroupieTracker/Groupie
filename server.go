package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

func HandleRegister(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return "err"
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
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
			return "err"
		}
		defer db.Close()

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM USER WHERE pseudo = ?", username)
		err = row.Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "err"
		}

		var count1 int
		row1 := db.QueryRow("SELECT COUNT(*) FROM USER WHERE email = ?", email)
		err = row1.Scan(&count1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "err"
		}

		if count1 > 0 {
			userError = "Cette adresse mail est déjà utilisé"
		} else if count > 0 {
			userError = "Ce nom d'utilisateur est déjà utilisé"
		} else {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "MDP pas hasher", http.StatusInternalServerError)
				return "err"
			}

			_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return "err"
			}

			http.Redirect(w, r, "/lobby", http.StatusSeeOther)
			return "err"
		}
	}
	registerError(w, userError)
	http.Redirect(w, r, "/lobby", http.StatusSeeOther)
	return username
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

func HandleLogin(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return "err"
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
	}

	usernameOrEmail := r.Form.Get("username")
	password := r.Form.Get("password")

	db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
	}
	defer db.Close()

	var storedPassword string
	row := db.QueryRow("SELECT password FROM USER WHERE pseudo = ? OR email = ?", usernameOrEmail, usernameOrEmail)
	err = row.Scan(&storedPassword)
	if err != nil {
		loginError(w, "Nom d'utilisateur ou mot de passe incorrect")
		return "err"
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		loginError(w, "Nom d'utilisateur ou mot de passe incorrect")
		return "err"
	}

	username, err := Groupi.GetUsernameByEmailOrUsername(db, usernameOrEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
	}
	expiration := time.Now().Add(24 * time.Hour)
	cookieName := "auth_token"
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    username,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/lobby", http.StatusSeeOther)

	return username
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

func GoBlindTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func GoGuessTheSong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/guessTheSong.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func GoScattergories(w http.ResponseWriter, r *http.Request, username string) {
	tmpl, err := template.ParseFiles("./static/scattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, username)
}

func ruleScattergories(r *http.Request) (string, int, int, int) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		name := r.FormValue("name")
		nbPlayer, _ := strconv.Atoi(r.FormValue("nbPlayer"))
		time, _ := strconv.Atoi(r.FormValue("time"))
		round, _ := strconv.Atoi(r.FormValue("nbRound"))
		return name, nbPlayer, time, round

	}
	return "", -1, -1, -1
}

func GoLobScattergories(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/lobbyScattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	var time int
	var nbRound int
	var username string
	http.HandleFunc("/", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/lobby", Lobby)
	http.HandleFunc("/BlindTest/webs", Groupi.WsBlindTest)
	http.HandleFunc("/GuessTheSong/webs", Groupi.WsGuessTheSong)
	http.HandleFunc("/LobScattergories", GoLobScattergories)
	http.HandleFunc("/logout", Logout)

	http.HandleFunc("/handle-login", func(w http.ResponseWriter, r *http.Request) {
		username = HandleLogin(w, r)
		if username == "err" {
			fmt.Println("err in login func")
		}

	})
	http.HandleFunc("/handle-register", func(w http.ResponseWriter, r *http.Request) {
		username = HandleRegister(w, r)
		if username == "err" {
			fmt.Println("err in login func")
		}

	})

	http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
		GoBlindTest(w, r)
	})
	http.HandleFunc("/GuessTheSong", func(w http.ResponseWriter, r *http.Request) {
		GoGuessTheSong(w, r)
	})
	http.HandleFunc("/Scattergories/webs", func(w http.ResponseWriter, r *http.Request) {
		Groupi.WsScattergories(w, r, time, nbRound, username)
	})

	http.HandleFunc("/RuleForScattergories", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
		defer db.Close()
		nameRooms, nbPlayer, ti, nbRo := ruleScattergories(r)
		time = ti
		nbRound = nbRo
		newGame := Groupi.Game{
			Name: nameRooms,
		}
		gameID, err := Groupi.CreateGameAndGetID(db, newGame)
		if err != nil {
			fmt.Println("Erreur lors de la création du jeu:", err)

			return
		}
		userID, err := Groupi.GetUserIDByUsername(db, username)
		if err != nil {
			fmt.Println("Erreur lors de la récupération de l'ID de l'utilisateur:", err)
			return
		}
		newRoom := Groupi.Roomms{
			CreatedBy:  userID,
			MaxPlayers: nbPlayer,
			Name:       nameRooms,
			GameID:     gameID,
		}
		roomID, err := Groupi.CreateRoomAndGetID(db, newRoom)
		id := strconv.Itoa(roomID)
		http.Redirect(w, r, "/Scattergories?room="+id, http.StatusSeeOther)

	})

	http.HandleFunc("/Scattergories", func(w http.ResponseWriter, r *http.Request) {
		GoScattergories(w, r, username)
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("static/", fs))

	fsCSS := http.FileServer(http.Dir("static/css"))
	http.Handle("/static/css/", http.StripPrefix("/static/css/", fsCSS))

	fsPicture := http.FileServer(http.Dir("static/assets/pictures"))
	http.Handle("/static/assets/pictures/", http.StripPrefix("/static/assets/pictures/", fsPicture))

	fsTracks := http.FileServer(http.Dir("static/assets/tracks"))
	http.Handle("/static/assets/tracks/", http.StripPrefix("/static/assets/tracks/", fsTracks))

	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
