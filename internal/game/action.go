package game

import "github.com/darkphotonKN/journey-through-midnight/internal/model"

/**
* Responsible for Game Events.
* Each round in the game, for each player, is an "Event".
* An event passes time and provides changes to the player's core state.
**/
type GameEvent struct {
	Type EventType
	Name string
}

type EventType string

const (
	Fight          EventType = "fight"
	PlayerOpponent EventType = "player_opponent"
	FoundItem      EventType = "found_item"
	Encounter      EventType = "encounter" // details TBD
	Shop           EventType = "shop"
)

/**
* Initiates a random event, waits for player choices, and outputs the results.
**/
// TODO: initiate a random event and play out the results
func (e *GameEvent) initiateEvent(eventType EventType) *model.PlayerState {
	return &model.PlayerState{}
}

/**
* Shop Event
* The player has the ability to buy one of many things offered from the shop.
* Provides a bunch of items, waits for response, adds chosen item to player's
* arsenal and deducts their gold.
**/
func (e *GameEvent) runShopEvent() *model.PlayerState {
	// send player list of purchaseable items

	// allow purchase

	// skip
	return nil
}
