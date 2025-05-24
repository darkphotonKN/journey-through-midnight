package matchmaking

import (
	"fmt"
	"sync"
	"time"

	"github.com/darkphotonKN/journey-through-midnight/internal/game"
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

type MatchMaker interface {
	JoinMatchMaking(*model.Player) error
	StartMatchMaking(interval time.Duration)
	GetQueueForTesting() []*model.Player
	GetNewGameChan() <-chan *game.Game
}

type BaseMatchMaker struct {
	queue []*model.Player
	mu    sync.Mutex

	// inject game engine to be able to initiate games directly via matchmaker
	game game.GameFactory

	// communication between server and matchmaking process to receive new games
	newGameChan chan *game.Game
}

// NOTE: neeed to inject GameFactory for matchmaker to be able to create its own game instances
func NewMatchMaker(gameFactory game.GameFactory) MatchMaker {
	return &BaseMatchMaker{
		queue:       make([]*model.Player, 0),
		game:        gameFactory,
		newGameChan: make(chan *game.Game),
	}
}

/**
* Allows a player to join matchmaking.
**/
func (m *BaseMatchMaker) JoinMatchMaking(player *model.Player) error {
	fmt.Printf("Player %s has queued for a game.\n", player.UserName)

	// TODO: can also add a map with id as key that updates at the same time
	// a user joins the queues to improve checking time complexity to O(1).
	for _, playerInQueue := range m.queue {
		if playerInQueue.ID == player.ID {
			fmt.Println("Player with the same ID attempted to join.")

			return fmt.Errorf("Player has already queued.")
		}
	}

	m.queue = append(m.queue, player)

	return nil
}

func (m *BaseMatchMaker) GetNewGameChan() <-chan *game.Game {
	return m.newGameChan
}

/**
* StartMatchmaking intializes the observation of players in the queue in a goroutine.
* This is currently started at the initialization of the game server.
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

	fmt.Printf("\nplayer one: %s, player two: %s\n\n", playerOne.UserName, playerTwo.UserName)

	players := []*model.Player{playerOne, playerTwo}
	newGameInstance := m.game.CreateGame(players)

	m.newGameChan <- newGameInstance

	// remove them from the queue
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
