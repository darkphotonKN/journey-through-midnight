package game

import "github.com/darkphotonKN/journey-through-midnight/internal/model"

/**
* Engine page are the core inner-workings of the game, holding together
* the logic, events, turns, etc.
**/

const (
	startingGold     int = 10
	maxInventorySize int = 5
)

func InitializeGameEngine() *GameFactory {

	initialConditions := InitialConditions{
		RoundDefault: 1,
		PlayerDefaults: PlayerDefaults{
			DefaultGold:  startingGold,
			DefaultItems: make([]model.Item, 5),
		},
	}

	// allow creating a game with factory
	defaultGameCreator := NewGameFactory(initialConditions)

	return defaultGameCreator
}
