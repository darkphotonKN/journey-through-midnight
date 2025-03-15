package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/darkphotonKN/journey-through-midnight/internal/game"
	"github.com/darkphotonKN/journey-through-midnight/internal/matchmaking"
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"

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

	// channel to send client message information to different goroutines
	serverChan chan ClientPackage

	// players concurrently online
	playersOnline map[uuid.UUID]model.Player

	// player connection to UUID mapping
	connToPlayerID map[*websocket.Conn]uuid.UUID

	// all current game connections
	games map[uuid.UUID]model.Game

	// stores unique ws connections for writing back to each client
	gameMsgChan map[*websocket.Conn]chan GameMessage

	// allows launching a game and matchmaking
	matchMaker matchmaking.MatchMaker

	// other
	mu sync.Mutex

	ListenAddr string

	upgrader websocket.Upgrader
}

const (
	matchmake_duration time.Duration = 15
)

func NewServer(listenAddr string) *Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Allow all connections by default for simplicity; can add more logic here
			return true
		},
	}

	// inject game factory for match maker to be able to create it's own game instances
	gameFactory := game.InitializeGameEngine()

	// instantiate a new matchmaker
	matchMaker := matchmaking.NewMatchMaker(*gameFactory)

	newServer := &Server{
		ListenAddr:     listenAddr,
		upgrader:       upgrader,
		serverChan:     make(chan ClientPackage),
		playersOnline:  make(map[uuid.UUID]model.Player),
		games:          make(map[uuid.UUID]model.Game),
		connToPlayerID: make(map[*websocket.Conn]uuid.UUID),
		gameMsgChan:    make(map[*websocket.Conn]chan GameMessage),
		matchMaker:     matchMaker,
	}

	// start matchmaking goroutine
	newServer.matchMaker.StartMatchMaking(time.Second * matchmake_duration)

	return newServer
}

func (s *Server) findPlayerByConnection(conn *websocket.Conn) (*model.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// find player id via their connection
	if id, exists := s.connToPlayerID[conn]; exists {

		// find corresponding player
		if player, ok := s.playersOnline[id]; ok {
			return &player, nil
		}
	}

	return nil, fmt.Errorf("Player with this connection does not exist.")
}
