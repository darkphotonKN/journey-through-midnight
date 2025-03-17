package game

import "github.com/darkphotonKN/journey-through-midnight/internal/model"

/**
* Engine page are the core inner-workings of the game, holding together
* the logic, events, turns, etc.
**/

type GameEngine struct {
	eventHandler *EventHandler
}

/**
* Handles the logic for ending the game round.
**/
func (e *GameEngine) ProcessEndRound(game *model.Game) {
	// increment round
	game.Round++
}

/**
* Process's a single player's round
**/
func (e *GameEngine) ProcessPlayerRound(game *model.Game, event model.GameEvent, playerState *model.PlayerState) {
	defaultEventHours := 1

	// play out event
	playerState = e.eventHandler.initiateEvent(event.Type)

	// increment their hours after event
	playerState.Time = model.Time{
		Day:  playerState.Time.Day,
		Hour: playerState.Time.Hour + defaultEventHours,
	}

	// update player state
	game.Players[playerState.ID] = playerState
}
