package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	Groupi"Groupi/Groupi"
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

	if len(password) < 8 || !containsDigit(password) || !containsLetter(password) || !containsSpecialChar(password) {
		data := struct {
			Error string
		}{
			Error: "Le mot de passe doit contenir au moins 8 caractères, inclure au moins un chiffre, une lettre et un caractère spécial",
		}
		tmpl, err := template.ParseFiles("static/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

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

	if count > 0 {
		http.Error(w, "Nom d'utilisateur déjà utilisé", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/lobby", http.StatusSeeOther)
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
	username := r.Form.Get("username")
	password := r.Form.Get("password")

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

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/lobby", http.StatusSeeOther)
}









































func GoBlindTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./pages/blindTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,nil)
}


func GOLobbyOfScattergories()  {
	//Creation d'une nouvelle party 
	//type petiti bac 
	//REcupere L'id de l'useur et le mais en t'en que createur


}

func GoGuessTheSong(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./pages/guessTheSong.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,nil)
}
func 	GoScattergories(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./pages/scattergories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,nil)
}



func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/lobby", Lobby)
	http.HandleFunc("/handle-register", HandleRegister)
	http.HandleFunc("/handle-login", HandleLogin)
// 
http.HandleFunc("/BlindTest/webs", Groupi.WsBlindTest)
http.HandleFunc("/GuessTheSong/webs", Groupi.WsGuessTheSong)
http.HandleFunc("/Scattergories/webs",Groupi.WsScattergories)
// URL pour le rendu de la page du blind test
http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
	GoBlindTest(w, r)
})
http.HandleFunc("/GuessTheSong", func(w http.ResponseWriter, r *http.Request) {
	GoGuessTheSong(w, r)
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