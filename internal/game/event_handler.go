package game

import (
	"math/rand"
	"time"
)

type Event string

const (
	eventStory Event = "storyEvent"
	eventShop  Event = "shopEvent"
)

/**
* Event handler processes events and their outcome.
**/
type EventHandler struct {
	events []Event
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

const (
	defaultEventHours = 1 // default increment of time when an event occurs
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

	// choose a random event
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randomIndex := rng.Intn(len(h.events))
	randomEvent := h.events[randomIndex]
}

/**
* Story Event
* This allows the player to pick choices based on the story and has a chance to result
* in earning an item or upgrade.
**/
type StoryEvent string

const (
	StoryEventPureStory StoryEvent = "pureStory"
	StoryEventEarnXP    StoryEvent = "earnXP"
	StoryEventEarnItem  StoryEvent = "earnItem"
)

var storyEvents = []StoryEvent{StoryEventPureStory, StoryEventEarnXP, StoryEventEarnItem}

func (h *EventHandler) runStoryEvent() StoryEvent {
	// choose a random event
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randomIndex := rng.Intn(len(storyEvents))
	selectedEvent := storyEvents[randomIndex]

	// selectedEvent now contains the randomly chosen story event
	return selectedEvent
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
