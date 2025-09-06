package server

import (
	"errors"
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

	// all current on-going game connections
	games map[uuid.UUID]*game.Game

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
	gameFactory := game.InitializeNewGameFactory()

	// instantiate a new matchmaker
	matchMaker := matchmaking.NewMatchMaker(*gameFactory)

	newServer := &Server{
		ListenAddr:     listenAddr,
		upgrader:       upgrader,
		serverChan:     make(chan ClientPackage),
		playersOnline:  make(map[uuid.UUID]model.Player),
		games:          make(map[uuid.UUID]*game.Game),
		connToPlayerID: make(map[*websocket.Conn]uuid.UUID),
		gameMsgChan:    make(map[*websocket.Conn]chan GameMessage),
		matchMaker:     matchMaker,
	}

	// start matchmaking goroutine
	newServer.matchMaker.StartMatchMaking(time.Second * matchmake_duration)

	return newServer
}

/**
* Finds a player by their unique websocket connection, in the case that
* id is unknown or not straight forward to acquire.
**/
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

/**
* Adds a created game to the list of games tracked on the server.
**/
func (s *Server) addGameToServer(g *game.Game) error {

	_, exists := s.games[g.ID]

	if exists {
		return game.ErrGameExists
	}

	// add to server's list of games
	s.games[g.ID] = g

	fmt.Printf("\ncurrent list of games: %+v\n\n", s.games)
	return nil
}

/**
* Find existing player in game instance.
**/
func (s *Server) findGameWithPlayer(id uuid.UUID) (*game.Game, error) {
	// loop through every single game and try to find if the player exists
	for _, game := range s.games {
		// find player inside game instance
		for _, player := range game.Players {
			if player.ID == id {
				return game, nil
			}
		}
	}

	return nil, errors.New("No player with this connection was found with any game.")
}

// NOTE: Methods for only testing
func (s *Server) GetMatchmaker() matchmaking.MatchMaker {
	return s.matchMaker
}
