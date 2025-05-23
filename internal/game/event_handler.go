package game

/**
* Event handler processes events and their outcome.
**/
type EventHandler struct {
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

const (
	defaultEventHours = 1
)

/**
* Initiates a random event, waits for player choices, and outputs the results.
**/
// TODO: initiate a random event and play out the results
func (h *EventHandler) initiateEvent(eventType EventType, playerState *PlayerState) {

	// increment their hours after event
	playerState.Time = Time{
		Day:  playerState.Time.Day,
		Hour: playerState.Time.Hour + defaultEventHours,
	}

}

/**
* Shop Event
* The player has the ability to buy one of many things offered from the shop.
* Provides a bunch of items, waits for response, adds chosen item to player's
* arsenal and deducts their gold.
**/
func (h *EventHandler) runShopEvent() *PlayerState {
	// send player list of purchaseable items

	// allow purchase

	// skip
	return nil
}
