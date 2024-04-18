package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	Groupi "Groupi/Groupi"
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

	if len(password) < 8 || !containsDigit(password) || !containsLetter(password) || !containsSpecialChar(password) {
		data := struct {
			Error string
		}{
			Error: "Le mot de passe doit contenir au moins 8 caractères, inclure au moins un chiffre, une lettre et un caractère spécial",
		}
		tmpl, err := template.ParseFiles("static/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "err"
		}
		tmpl.Execute(w, data)
		return "err"
	}

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
		return"err"
	}

	if count > 0 {
		http.Error(w, "Nom d'utilisateur déjà utilisé", http.StatusBadRequest)
		return "err"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
		return "err"
	}

	_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
	}

	http.Redirect(w, r, "/lobby", http.StatusSeeOther)
	return username
}

func containsDigit(s string) bool {
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
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
	}
	defer db.Close()

	var storedPassword string
	row := db.QueryRow("SELECT password FROM USER WHERE pseudo = ?", username)
	err = row.Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
		return "err"
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
		return "err"
	}

	http.Redirect(w, r, "/lobby", http.StatusSeeOther)
	return username
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
func GoScattergories(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/scattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ruleScattergories( r *http.Request) (string,int,int,int) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		name := r.FormValue("name")
		fmt.Println(r.FormValue("nbPlayer"))
		nbPlayer,_ := strconv.Atoi( r.FormValue("nbPlayer"))
		time ,_:=  strconv.Atoi(r.FormValue("time"))
		round,_ := strconv.Atoi( r.FormValue("nbRound"))
		fmt.Println(name, nbPlayer, time, round)
		return name, nbPlayer ,time , round
		
	}
	return "",-1,-1,-1
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
	http.HandleFunc("/GuessTheSong/webs", Groupi.WsGuessTheSong )
	http.HandleFunc("/LobScattergories", GoLobScattergories)

	http.HandleFunc("/handle-login", func(w http.ResponseWriter, r *http.Request) {
		username=HandleLogin(w, r)
		if username == "err"{
			fmt.Println("err in login func")
		}
		fmt.Println(username)
	})
	http.HandleFunc("/handle-register", func(w http.ResponseWriter, r *http.Request) {
		username=HandleRegister(w, r)
		if username == "err"{
			fmt.Println("err in login func")
		}
		fmt.Println(username)
	})

	
	http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
		GoBlindTest(w, r)
	})
	http.HandleFunc("/GuessTheSong", func(w http.ResponseWriter, r *http.Request) {
		GoGuessTheSong(w, r)
	})
	http.HandleFunc("/Scattergories/webs", func(w http.ResponseWriter, r *http.Request) {
		Groupi.WsScattergories(w, r,time,nbRound, username)
	})

	http.HandleFunc("/RuleForScattergories", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
		defer db.Close()
		nameRooms,nbPlayer,ti,nbRo:=ruleScattergories(r)
		time=ti
		nbRound=nbRo
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
		GoScattergories(w, r)
	})
	//

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fsCSS := http.FileServer(http.Dir("static/css"))
	http.Handle("/static/css/", http.StripPrefix("/static/css/", fsCSS))

	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
