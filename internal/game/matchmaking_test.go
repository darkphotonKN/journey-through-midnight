package game

import (
	"time"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

func TestMatchMaker_MatchMake() {

	matchMaker := NewMatchMaker()

	matchMaker.StartMatchMaking(time.Second * 5)

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

	// assert for only 1 player left after 1 second

	// TODO: assert for players 1 and 2 not in queue
}
