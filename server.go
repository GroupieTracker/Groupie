package main 

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	Groupi"Groupi/Groupi"
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
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        Home(w, r)
    })

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

	
    http.HandleFunc("/BlindTest/webs", Groupi.WsBlindTest)
	http.HandleFunc("/GuessTheSong/webs", Groupi.WsGuessTheSong)
	http.HandleFunc("/Scattergories/webs",Groupi.WsScattergories)



    // Serveur de fichiers statiques
    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    fmt.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
