package game

import (
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

/**
* Factory for creating games with default settings and player state
* conversion.
**/
type GameFactory struct {
	defaultRound       int
	defaultPlayerState DefaultPlayerState
}

type PlayerDefaults struct {
	DefaultGold  int
	DefaultItems []Item
}

type InitialConditions struct {
	RoundDefault   int
	PlayerDefaults PlayerDefaults
}

type DefaultPlayerState struct {
	defaultGold       int
	defaultItems      []Item
	defaultAttributes Attributes
}

func NewGameFactory(initialConditions InitialConditions) *GameFactory {

	return &GameFactory{
		defaultRound: initialConditions.RoundDefault,
		defaultPlayerState: DefaultPlayerState{
			defaultGold:  initialConditions.PlayerDefaults.DefaultGold,
			defaultItems: initialConditions.PlayerDefaults.DefaultItems,
			defaultAttributes: Attributes{
				Strength:     1,
				Intelligence: 1,
				Wisdom:       1,
				Agility:      1,
				Vitality:     1,
				Faith:        1,
				Charisma:     1,
			},
		},
	}
}

/**
* Core game creation (constuctor), based on factory settings and converting all passed
* in players into "PlayerState"s.
**/
func (f *GameFactory) CreateGame(players []*model.Player) *Game {
	// convert Players to PlayerState
	playerStates := make(map[uuid.UUID]*PlayerState)

	for _, player := range players {
		playerStates[player.ID] = &PlayerState{
			Player:    *player,
			Gold:      f.defaultPlayerState.defaultGold,
			Inventory: f.defaultPlayerState.defaultItems,
		}
	}

	eventHandler := NewEventHandler()

	return &Game{
		ID:      uuid.New(),
		MsgChan: make(chan string),
		Round:   f.defaultRound,
		Players: playerStates,

		// inject eventHandler
		eventHandler: eventHandler,
	}
}

const (
	startingGold     int = 10
	maxInventorySize int = 5
)

func InitializeNewGameFactory() *GameFactory {

	initialConditions := InitialConditions{
		RoundDefault: 1,
		PlayerDefaults: PlayerDefaults{
			DefaultGold:  startingGold,
			DefaultItems: make([]Item, 5),
		},
	}

	// allow creating a game with factory
	defaultGameCreator := NewGameFactory(initialConditions)

	return defaultGameCreator
}
