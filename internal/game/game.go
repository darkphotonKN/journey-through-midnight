package game

import "github.com/darkphotonKN/journey-through-midnight/internal/model"

import (
	"github.com/google/uuid"
)

/**
* Holds all the information for a specific game's meta data.
**/
type Game struct {
	MsgChan      GameMsgChan   // message channel to communicate with game
	eventHandler *EventHandler // handles game events

	// --- metadata ---

	ID    uuid.UUID // unique identifier for each game
	Round int       // also represents "day"

	// players in this game instance
	Players map[uuid.UUID]*PlayerState
}

type GameMsgChan = chan string

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

type Time struct {
	Day  int
	Hour int
}

type HeroClassName string

const (
	fighter HeroClassName = "Fighter"
	wizard  HeroClassName = "Wizard"
	rogue   HeroClassName = "Rogue"
	priest  HeroClassName = "Priest"
	duelist HeroClassName = "Duelist"
	templar HeroClassName = "Templar"
)

type HeroClass struct {
}

type Attributes struct {
	Strength     int
	Intelligence int
	Wisdom       int
	Agility      int
	Vitality     int
	Faith        int
	Charisma     int
}

type Item struct {
	id   uuid.UUID
	name string
}

type Skill struct {
	id   uuid.UUID
	name string
}

/**
* Holds all primary information about a player.
**/
type PlayerState struct {
	model.Player
	// the stage the game match has reached for any player time Time
	HeroInfo  Hero       // hero information
	Party     []Follower // TODO: update to include followers
	Inventory []Item     // global items
	Gold      int
	Time      Time // time the player has reached
}

type Hero struct {
	class      HeroClass  // mage
	attributes Attributes //
	skills     []Skill
	items      []Item // items specifically meant for a hero
}

type Follower struct {
	heroInfo Hero
}
