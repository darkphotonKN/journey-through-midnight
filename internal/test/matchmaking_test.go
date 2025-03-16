package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/darkphotonKN/journey-through-midnight/internal/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMatchMaker_MatchMake(t *testing.T) {

	// simulate game server and message hub
	s := server.NewServer("7777")
	go s.MessageHub()

	matchMaker := s.GetMatchmaker()

	matchWaitTime := time.Second * 1
	matchMaker.StartMatchMaking(matchWaitTime) // overwrite original game wait time

	playerOneUUID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	playerTwoUUID, _ := uuid.Parse("11111111-1111-1111-1111-111111111112")
	playerThreeUUID, _ := uuid.Parse("11111111-1111-1111-1111-111111111113")

	// add three players to queue
	playerOne := model.Player{
		ID:       playerOneUUID,
		UserName: "Player 1: John",
	}

	playerTwo := model.Player{
		ID:       playerTwoUUID,
		UserName: "Player 2: Sandra",
	}

	playerThree := model.Player{
		ID:       playerThreeUUID,
		UserName: "Player 3: Bullock",
	}

	matchMaker.JoinMatchMaking(&playerOne)
	matchMaker.JoinMatchMaking(&playerTwo)
	matchMaker.JoinMatchMaking(&playerThree)

	queue := matchMaker.GetQueueForTesting()

	fmt.Printf("\ncurrent queue: %+v\n\n", queue)

	// test for length of initial queue
	assert.Len(t, queue, 3)

	// assert for only 1 player left after wait time
	waitTimeOffset := time.Millisecond * 500
	timer := time.NewTicker(matchWaitTime + waitTimeOffset)

	select {
	case <-timer.C:
		queue = matchMaker.GetQueueForTesting()

		assert.Len(t, queue, 1)
	}

	// TODO: assert for players 1 and 2 not in queue
}
