package game

import (
	"fmt"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

type MatchMaker interface {
	InitiateMatchMaking(*model.Player) error
	TryCreateMatch() (bool, uuid.UUID, model.GameInformation)
}

type BaseMatchMaker struct {
	// players in queue
	queue []*model.Player
}

func NewMatchMaker() MatchMaker {
	return &BaseMatchMaker{
		queue: make([]*model.Player, 0),
	}
}

/**
* Initiates match making to put players into a game.
* TODO: Right now only adds player to player online list, no real match
* making (in order to develop core game logic first).
**/
func (m *BaseMatchMaker) InitiateMatchMaking(player *model.Player) error {

	fmt.Printf("Player %s has queued for a game.\n", player.UserName)
	// adds player to queue
	m.queue = append(m.queue, player)

	return nil
}

/**
* Attempts to create a match with currently queued players.
**/
func (m *BaseMatchMaker) TryCreateMatch() (bool, uuid.UUID, model.GameInformation) {

	return false, uuid.New(), model.GameInformation{}
}
