package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

type MatchMaker interface {
	JoinMatchMaking(*model.Player) error
	StartMatchMaking(interval time.Duration)
	TryCreateMatch() (bool, uuid.UUID, model.Game)
	GetQueueForTesting() []*model.Player
}

type BaseMatchMaker struct {
	// players in queue
	queue []*model.Player

	// matchmaker's own mutex
	mu sync.Mutex
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
	// NOTE: only debug logging
	var playersInQueue []model.Player
	for _, player := range m.queue {
		playersInQueue = append(playersInQueue, *player)
	}
	fmt.Printf("\nInitial queue: %+v\n\n", playersInQueue)
	// NOTE: debug logging end

	if len(m.queue) < 2 {
		fmt.Printf("\nNot enough players in queue, waiting on players...\n\n")
		return
	}

	// check queue and take first 2 players into the match
	playerOne := m.queue[0]
	playerTwo := m.queue[1]

	if playerOne == nil || playerTwo == nil {
		fmt.Printf("Error when matchmaking, required players don't exist in the queue\n")
		return
	}

	// TODO: creates game with these players
	fmt.Printf("\nplayer one: %s, player two: %s\n\n", playerOne.UserName, playerTwo.UserName)

	// remove them from queue
	m.removePlayerFromQueue(playerOne.ID)
	m.removePlayerFromQueue(playerTwo.ID)

	// NOTE: only debug logging
	playersInQueue = []model.Player{}
	for _, player := range m.queue {
		playersInQueue = append(playersInQueue, *player)
	}
	fmt.Printf("\nRemaining queue: %+v\n\n", playersInQueue)
	// NOTE: debug logging end
}

/**
* Removes players from the queue.
**/
func (m *BaseMatchMaker) removePlayerFromQueue(id uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for index, player := range m.queue {

		if player.ID == id {
			m.removeByIndex(index)
		}
	}
}

/**
* Remove player from queue by index (memory efficient version)
**/
func (m *BaseMatchMaker) removeByIndex(index int) {
	// replace to last index
	m.queue[index] = m.queue[len(m.queue)-1]

	// truncate away the element
	m.queue = m.queue[:len(m.queue)-1]
}

// --- NOTE: METHODS ONLY FOR TESTING ---
func (m *BaseMatchMaker) GetQueueForTesting() []*model.Player {
	return m.queue
}
