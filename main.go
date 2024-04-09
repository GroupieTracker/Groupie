package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var (
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
    connections = make(map[*websocket.Conn]bool)
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    connections[conn] = true

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            delete(connections, conn)
            return
        }

        fmt.Println(string(p))

        for conn := range connections {
            if err := conn.WriteMessage(messageType, p); err != nil {
                log.Println(err)
                conn.Close()
                delete(connections, conn)
            }
        }
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func main() {
    http.HandleFunc("/echo", echoHandler)
    http.HandleFunc("/", homeHandler)

    fmt.Println("Serveur WebSocket démarré sur le port 8080...")
    if err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil); err != nil {
        log.Fatal("Serveur error:", err)
    }
}
