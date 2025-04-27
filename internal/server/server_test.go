package server

import (
	"testing"

	"github.com/darkphotonKN/journey-through-midnight/internal/game"
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	// setup test values to be shared across test methods
	server    *Server
	playerOne model.Player
	playerTwo model.Player
}

func (s *ServerTestSuite) SetupTest() {
	newServer := NewServer(testAddr)

	s.server = newServer

	// setup dummy clients
	s.playerOne = model.Player{
		ID:       uuid.New(),
		UserName: "Player One",
	}

	s.playerTwo = model.Player{
		ID:       uuid.New(),
		UserName: "Player Two",
	}

	// setup dummy games
	factory := game.InitializeNewGameFactory()

	// add dummy players to dummy game
	var testPlayers []*model.Player

	testPlayers = append(testPlayers, &s.playerOne)

	testPlayers = append(testPlayers, &s.playerTwo)

	testGame := factory.CreateGame(testPlayers)

	s.server.addGameToServer(testGame)
}

// test games can be successfully added
func (s *ServerTestSuite) TestGameAddedToServer() {
	assert.Equal(s.T(), 1, len(s.server.games))
}

// test games added players can be found in the games
func (s *ServerTestSuite) TestPlayerFindable() {

	game, err := s.server.findGameWithPlayer(s.playerOne.ID)
	assert.NoError(s.T(), err)

	// game must exist
	assert.NotNil(s.T(), game)

	// game with non-existant player should not exist
	gameTwo, _ := s.server.findGameWithPlayer(uuid.New())
	assert.Nil(s.T(), gameTwo)
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
