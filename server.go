package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	Groupi "Groupi/Groupi"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Redirect(w, r, "/lobby", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "./static/index.html")
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
	tmpl, err := template.ParseFiles("./static/scattergories/scattergories.html")
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
	tmpl, err := template.ParseFiles("./static/scattergories/lobbyScattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoListScattergories(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/scattergories/listeScattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("sqlite3", "./Groupi/BDD.db")
		defer db.Close()
	tab , _ :=Groupi.GetRoomsByGameCategory(db , "scattergories")
	tmpl.Execute(w, tab)
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
	http.HandleFunc("/ListLobOfScattergories", GoListScattergories)
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
			fmt.Println("err in register func")
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
		// Groupi.ClearDatabase(db)
		defer db.Close()
		nameRooms, nbPlayer, ti, nbRo := ruleScattergories(r)
		time = ti
		nbRound = nbRo
		newGame := Groupi.Game{
			Name: "scattergories",
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
