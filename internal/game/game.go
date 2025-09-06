package game

import (
	"sync"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

/**
* Holds all the information for a specific game's meta data.
**/
type Game struct {
	MsgCh       GameMsgChan // message channel to communicate actions
	CloseGameCh chan bool   // for communicating the closing of an existing game

	eventHandler *EventHandler // handles game events

	// --- metadata ---

	ID    uuid.UUID // unique identifier for each game
	Round int       // NOTE: "Round" also represents the "day" of the game
	Phase Phase

	// players in this game instance
	Players map[uuid.UUID]*PlayerState

	// mutex
	mu sync.Mutex
}

type GameMsgChan = chan interface{}

/**
* Responsible for Game Events.
* Each round in the game, for each player, is an "Event".
* An event passes time and provides changes to the player's core state.
**/
type GameEvent struct {
	Type EventType
	Name string
}

type Phase string

const (
	Day      Phase = "Day"
	Midnight Phase = "Midnight"
)

type EventType string

const (
	Fight          EventType = "fight"
	PlayerOpponent EventType = "player_opponent"
	FoundItem      EventType = "found_item"
	Encounter      EventType = "encounter"
	Shop           EventType = "shop"
)

type Time struct {
	Day  int
	Hour int
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
	class      HeroClass
	attributes Attributes
	skills     []Skill
	items      []Item
	stats      Stats
}

type Stats struct {
	health    int // health for each fight
	endurance int // total health for the entire night
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
	name  HeroClassName
	level int
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

type Follower struct {
	heroInfo Hero
}
