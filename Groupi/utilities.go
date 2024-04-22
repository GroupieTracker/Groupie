package Groupi

import (
	"net/http"
	"sync"

	websocket "github.com/gorilla/websocket"
)

type Room struct {
	ID          string
	Connections map[*websocket.Conn]bool
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rooms = make(map[string]*Room)
	mutex = sync.Mutex{}
)
