package test

import (
	"testing"
	"time"

	"github.com/darkphotonKN/journey-through-midnight/internal/matchmaking"
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/darkphotonKN/journey-through-midnight/internal/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MatchmakingTestSuite struct {
	suite.Suite

	// shared test values across test methods
	matchmaker    matchmaking.MatchMaker
	players       []*model.Player
	matchWaitTime time.Duration
}

func (s *MatchmakingTestSuite) SetupTest() {
	// simulate game server and message hub
	serv := server.NewServer("7777")
	go serv.MessageHub()

	matchMaker := serv.GetMatchmaker()

	matchWaitTime := time.Second * 1
	s.matchWaitTime = matchWaitTime
	matchMaker.StartMatchMaking(s.matchWaitTime) // overwrite original game wait time

	playerOneUUID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	playerTwoUUID, _ := uuid.Parse("11111111-1111-1111-1111-111111111112")
	playerThreeUUID, _ := uuid.Parse("11111111-1111-1111-1111-111111111113")

	// add three players to queue
	playerOne := &model.Player{
		ID:       playerOneUUID,
		UserName: "Player 1: John",
	}

	playerTwo := &model.Player{
		ID:       playerTwoUUID,
		UserName: "Player 2: Sandra",
	}

	playerThree := &model.Player{
		ID:       playerThreeUUID,
		UserName: "Player 3: Person",
	}

	s.players = make([]*model.Player, 0)

	s.players = append(s.players, playerOne)
	s.players = append(s.players, playerTwo)
	s.players = append(s.players, playerThree)

	s.matchmaker = matchMaker
}

func (s *MatchmakingTestSuite) TestQueue() {
	s.matchmaker.JoinMatchMaking(s.players[0])
	s.matchmaker.JoinMatchMaking(s.players[1])
	s.matchmaker.JoinMatchMaking(s.players[2])

	queue := s.matchmaker.GetQueueForTesting()

	playerThreeCount := 0

	for _, player := range queue {

		if player.ID == s.players[2].ID {
			playerThreeCount++
		}
	}

	// make sure player 3 only appears once despite joining twice
	assert.Equal(s.T(), playerThreeCount, 1)

	// test for length of initial queue
	assert.Len(s.T(), queue, 3)

	// assert for only 1 player left after wait time
	waitTimeOffset := time.Millisecond * 500
	timer := time.NewTicker(s.matchWaitTime + waitTimeOffset)

	select {
	case <-timer.C:
		queue = s.matchmaker.GetQueueForTesting()

		assert.Len(s.T(), queue, 1)
	}
}

func TestServer(t *testing.T) {
	suite.Run(t, new(MatchmakingTestSuite))
}
