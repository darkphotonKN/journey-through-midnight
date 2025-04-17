package server

import (
	"fmt"
	"testing"

	"github.com/darkphotonKN/journey-through-midnight/internal/game"
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

/**
* All tests related to server connections and goroutines.
**/

var (
	testAddr = ":5555"
)

type ServerTestSuite struct {
	suite.Suite
	server *Server
}

func (s *ServerTestSuite) SetupTest() {
	newServer := NewServer(testAddr)
	s.server = newServer
}

func (s *ServerTestSuite) TestClientCleanUp() {
	// setup dummy clients
	playerOne := model.Player{
		ID:       uuid.New(),
		UserName: "Player One",
	}

	playerTwo := model.Player{
		ID:       uuid.New(),
		UserName: "Player Two",
	}

	// setup dummy games
	testGame := game.Game{
		ID: uuid.New(),
	}

	// add dummy players to dummy game

	testGame.Players[playerOne.ID] = &game.PlayerState{
		Player: playerOne,
	}

	testGame.Players[playerTwo.ID] = &game.PlayerState{
		Player: playerTwo,
	}

	s.server.games[testGame.ID] = &testGame

	fmt.Printf("Game before client cleanup: \n%v\n\n", s.server.games[testGame.ID])

	fmt.Printf("Game after client cleanup: \n%v\n\n", s.server.games[testGame.ID])
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
