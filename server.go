package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
)

var messages []string

func main() {
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/send", handleSend)
    
    log.Println("Server started. Listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    tpl := template.Must(template.ParseFiles("index.html"))
    tpl.Execute(w, messages)
}

func handleSend(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    message := r.FormValue("message")
    messages = append(messages, message)
    fmt.Println("Message received:", message)
    
    // Redirect to the home page to display the updated messages
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
