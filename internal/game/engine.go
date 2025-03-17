package game

/**
* Engine page are the core inner-workings of the game, holding together
* the logic, events, turns, etc.
**/

/**
* Handles the logic for ending the game round.
**/
func (g *Game) ProcessEndRound() {
	// increment round
	g.Round++
}

/**
* Process's a single player's round
**/
func (g *Game) ProcessPlayerRound(game *Game, event GameEvent, playerState *PlayerState) {
	defaultEventHours := 1

	// play out event
	playerState = g.eventHandler.initiateEvent(event.Type)

	// increment their hours after event
	playerState.Time = Time{
		Day:  playerState.Time.Day,
		Hour: playerState.Time.Hour + defaultEventHours,
	}

	// update player state
	game.Players[playerState.ID] = playerState
}
