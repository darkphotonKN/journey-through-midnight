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

/**
* All the logic that surrounds starting an actual game.
**/

/**
* Handles the logic for ending the game round.
**/
func (g *model.Game) ProcessEndRound() {
	// increment round
	g.Round = g.Round + 1
}

/**
* Process's a single player's round
**/
func (g *model.Game) ProcessPlayerRound(event GameEvent, playerState *model.PlayerState) {
	defaultEventHours := 1

	// play out event
	playerState = event.initiateEvent(event.Type)

	// increment their hours after event
	playerState.Time = model.Time{
		Day:  playerState.Time.Day,
		Hour: playerState.Time.Hour + defaultEventHours,
	}

	// update player state
	g.Players[playerState.ID] = playerState
}
