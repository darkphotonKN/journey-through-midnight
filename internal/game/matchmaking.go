package game

import (
	"fmt"
	"time"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

type MatchMaker interface {
	JoinMatchMaking(*model.Player) error
	StartMatchMaking(interval time.Duration)
	TryCreateMatch() (bool, uuid.UUID, model.Game)
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
* Allows a player to matchmaking.
**/
func (m *BaseMatchMaker) JoinMatchMaking(player *model.Player) error {

	fmt.Printf("Player %s has queued for a game.\n", player.UserName)

	// adds player to queue
	m.queue = append(m.queue, player)

	return nil
}

/**
* Attempts to create a match with currently queued players.
**/
func (m *BaseMatchMaker) TryCreateMatch() (bool, uuid.UUID, model.Game) {

	return false, uuid.New(), model.Game{}
}

/**
* StartMatchmaking intializes the observation of players in the queue in a goroutine.
**/
func (m *BaseMatchMaker) StartMatchMaking(interval time.Duration) {

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			// channel delivering ticks
			case <-ticker.C:
				fmt.Println("one tick of matchmaking...")
				m.matchMake()
			}
		}

	}()
}

/**
* Checks queue and matchmakes players.
**/

func (m *BaseMatchMaker) matchMake() {

	if len(m.queue) < 2 {
		fmt.Printf("Not enough players in queue, waiting on players...\n")
		return
	}

	// check queue and take first 2 players into the match
	playerOne := m.queue[0]
	playerTwo := m.queue[0]

	if playerOne == nil || playerTwo == nil {
		fmt.Printf("Error when matchmaking, required players don't exist in the queue\n")
		return
	}

	// creates game with these players

	// remove them from queue
	m.queue = m.queue[2:]

	fmt.Printf("\nRemaining queue:\n\n", m.queue)
}
