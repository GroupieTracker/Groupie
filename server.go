package main 

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	
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
	tmpl.Execute(w, nil)
}


func main()  {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r)
	})
	http.HandleFunc("/goBlindTest", func(w http.ResponseWriter, r *http.Request) {
		GoBlindTest(w, r)
	})
	http.HandleFunc("/blindTest", func(w http.ResponseWriter, r *http.Request) {
		BlindTest(w, r)
	})
}

fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))