package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

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

func Team(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/team.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoBlindTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindtest/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoWinnerTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindtest/winner.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoGuessTheSong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/guessthesong/guessTheSong.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoScattergories(w http.ResponseWriter, r *http.Request, category []string) {
	tmpl, err := template.ParseFiles("./static/scattergories/scattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, category)
}

func GoWattingForScattergories(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/scattergories/roomWaiting.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ruleScattergories(r *http.Request) (string, int, int, int, []string) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		name := r.FormValue("name")
		nbPlayer, _ := strconv.Atoi(r.FormValue("nbPlayer"))
		time, _ := strconv.Atoi(r.FormValue("time"))
		round, _ := strconv.Atoi(r.FormValue("nbRound"))
		category := r.Form["categ"]

		return name, nbPlayer, time, round, category

	}
	return "", -1, -1, -1, nil
}

func ruleBlindtest(r *http.Request) (string, int, int, int) {
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

func GoListBlindtest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindtest/listeBlindtest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
	defer db.Close()
	tab, _ := Groupi.GetRoomsByGameCategory(db, "blindtest")
	tmpl.Execute(w, tab)
}

func GoLobBlindtest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindtest/lobblindtest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoLobGuessthesong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/guessthesong/lobguessthesong.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoListGuessthesong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/guessthesong/listeGuessthesong.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
	defer db.Close()
	tab, _ := Groupi.GetRoomsByGameCategory(db, "guessthesong")
	tmpl.Execute(w, tab)
}

func GoLobScattergories(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/scattergories/lobbyScattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GoWinner(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindtest/winner.html")
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
	db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
	defer db.Close()
	tab, _ := Groupi.GetRoomsByGameCategory(db, "scattergories")
	tmpl.Execute(w, tab)
}

func GoWattingForBlindTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/blindtest/waitingRoom.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func GoResult(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
	defer db.Close()
	queryParams := r.URL.Query()
	roomIDStr := queryParams.Get("room")
	id, _ := strconv.Atoi(roomIDStr)
	userScoreTab, _ := Groupi.GetUserScoresForRoomID(db, id)
	sort.Slice(userScoreTab, func(i, j int) bool {
		return userScoreTab[i][1] > userScoreTab[j][1]
	})

	tmpl.Execute(w, userScoreTab)
	time.Sleep(1 * time.Second)
	Groupi.DeleteRoomData(db, id)
}
func GoWattingForGuessTheSong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/guessthesong/waitingRoom.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	var time int
	var timeForRoundBlindTest int
	var nbRound int
	var nbRoundBlindTest int
	var username string
	var category []string
	var dif string
	db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
	defer db.Close()
	Groupi.ClearDatabase(db)
	http.HandleFunc("/", Home)
	http.HandleFunc("/team", Team)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/lobby", Lobby)
	http.HandleFunc("/GuessTheSong/webs", func(w http.ResponseWriter, r *http.Request) {

		Groupi.WsGuessTheSong(w, r, time, nbRound, dif)
	})
	http.HandleFunc("/LobBlindtest", GoLobBlindtest)
	http.HandleFunc("/ListBlindtest", GoListBlindtest)
	http.HandleFunc("/LobGuessthesong", GoLobGuessthesong)
	http.HandleFunc("/ListGuessthesong", GoListGuessthesong)
	http.HandleFunc("/WaitingRoomForBlindTest", GoWattingForBlindTest)
	http.HandleFunc("/WaitingRoomForGuessTheSong", GoWattingForGuessTheSong)
	http.HandleFunc("/WaitingRoomForScattergories/webs", Groupi.WsWaitingRoom)
	http.HandleFunc("/WaitingRoomForBlindTest/webs", Groupi.WsWaitingRoomBlindtest)
	http.HandleFunc("/WsWaitingRoomGuessTheSong/webs", Groupi.WsWaitingRoom)
	http.HandleFunc("/Result", GoResult)
	http.HandleFunc("/WaitingRoomForBlintest/webs", Groupi.WsWaitingRoomBlindtest)
	http.HandleFunc("/Winner/webs", GoWinner)

	http.HandleFunc("/LobScattergories", GoLobScattergories)
	http.HandleFunc("/ListLobOfScattergories", GoListScattergories)
	http.HandleFunc("/logout", Logout)

	http.HandleFunc("/handle-login", func(w http.ResponseWriter, r *http.Request) {
		HandleLogin(w, r)
	})
	http.HandleFunc("/handle-register", func(w http.ResponseWriter, r *http.Request) {
		HandleRegister(w, r)
	})

	http.HandleFunc("/BlindTest/webs", func(w http.ResponseWriter, r *http.Request) {
		Groupi.WsBlindTest(w, r, timeForRoundBlindTest, nbRoundBlindTest)
	})
	http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
		GoBlindTest(w, r)
	})

	http.HandleFunc("/LobblindTest", func(w http.ResponseWriter, r *http.Request) {
		GoLobBlindtest(w, r)
	})

	http.HandleFunc("ListBlindtest/winner", func(w http.ResponseWriter, r *http.Request) {
		GoWinnerTest(w, r)
	})

	http.HandleFunc("/GuessTheSong", func(w http.ResponseWriter, r *http.Request) {
		GoGuessTheSong(w, r)
	})
	http.HandleFunc("/Scattergories/webs", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			fmt.Println("Erreur lors de la récupération du cookie :", err)
			return
		}
		username = cookie.Value
		db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
		defer db.Close()
		roomID := r.URL.Query().Get("room")
		roomIDInt, _ := strconv.Atoi(roomID)
		userID, _ := Groupi.GetUserIDByUsername(db, username)
		Groupi.AddRoomUser(db, roomIDInt, userID)
		fmt.Println(roomIDInt, userID)
		usersIDs, _ := Groupi.GetUsersInRoom(db, roomID)
		fmt.Println(usersIDs)
		Groupi.WsScattergories(w, r, time, nbRound)
	})

	http.HandleFunc("/RuleForScattergories", func(w http.ResponseWriter, r *http.Request) {
		db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			fmt.Println("Erreur lors de la récupération du cookie :", err)
			return
		}
		username = cookie.Value

		defer db.Close()
		nameRooms, nbPlayer, ti, nbRo, cat := ruleScattergories(r)
		time = ti
		category = cat
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
			Category:   category,
		}
		roomID, _ := Groupi.CreateRoomAndGetID(db, newRoom)
		id := strconv.Itoa(roomID)
		http.Redirect(w, r, "/WaitingRoomForScattergories?room="+id, http.StatusSeeOther)

	})

	http.HandleFunc("/RuleForBlindtest", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			fmt.Println("Erreur lors de la récupération du cookie :", err)
			return
		}
		username = cookie.Value

		db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
		defer db.Close()
		nameRooms, nbPlayer, ti, nbRo := ruleBlindtest(r)
		timeForRoundBlindTest = ti
		nbRoundBlindTest = nbRo
		newGame := Groupi.Game{
			Name: "blindtest",
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
		roomID, _ := Groupi.CreateRoomAndGetID(db, newRoom)
		id := strconv.Itoa(roomID)
		http.Redirect(w, r, "/WaitingRoomForBlindTest?room="+id, http.StatusSeeOther)

	})

	http.HandleFunc("/RuleForGuessthesong", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			fmt.Println("Erreur lors de la récupération du cookie :", err)
			return
		}
		username = cookie.Value
		db, _ := sql.Open("sqlite3", "./Groupi/BDD.db")
		defer db.Close()
		nameRooms, nbPlayer, ti, nbRo := ruleBlindtest(r)
		time = ti
		nbRound = nbRo
		newGame := Groupi.Game{
			Name: "guessthesong",
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
		roomID, _ := Groupi.CreateRoomAndGetID(db, newRoom)
		id := strconv.Itoa(roomID)
		http.Redirect(w, r, "/WaitingRoomForGuessTheSong?room="+id, http.StatusSeeOther)

	})

	http.HandleFunc("/Scattergories", func(w http.ResponseWriter, r *http.Request) {
		GoScattergories(w, r, category)
	})

	http.HandleFunc("/WaitingRoomForScattergories", func(w http.ResponseWriter, r *http.Request) {
		GoWattingForScattergories(w, r)
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
