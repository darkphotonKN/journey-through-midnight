package server

import (
	"sync"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"

	"github.com/gorilla/websocket"

	"net/http"
)

// to group GameMessage along with the websocket connection to pass to message hub for handling
type ClientPackage struct {
	GameMessage GameMessage
	Conn        *websocket.Conn
}

// Core Server Type
type Server struct {
	// address to connect to the server
	listenAddr string

	// for upgrading connection to websocket
	upgrader websocket.Upgrader

	// players concurrently online
	playersOnline map[*websocket.Conn]model.Player

	// all current game connections
	games map[*websocket.Conn]model.GameInformation

	// other
	mu sync.Mutex
}

func NewServer(listenAddr string) *Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Allow all connections by default for simplicity; can add more logic here
			return true
		},
	}

	return &Server{
		listenAddr:    listenAddr,
		upgrader:      upgrader,
		playersOnline: make(map[*websocket.Conn]model.Player),
		games:         make(map[*websocket.Conn]model.GameInformation),
	}
}
