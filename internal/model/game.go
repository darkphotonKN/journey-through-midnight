package model

import "github.com/google/uuid"

type GameMsgChan = chan string

type Time struct {
	day  int
	hour int
}

type HeroClassName string

const (
	knight  HeroClassName = "knight"
	wizard  HeroClassName = "wizard"
	hunter  HeroClassName = "Hunter"
	monk    HeroClassName = "Monk"
	priest  HeroClassName = "Priest"
	duelist HeroClassName = "Duelist"
)

type HeroClass struct {
}

type Item struct {
	id   uuid.UUID
	name string
}

type Items []Item

type PlayerState struct {
	Player
	// the stage the game match has reached for any player
	time  map[uuid.UUID]Time
	class HeroClass
}

/**
* Holds all the information for a specific game's meta data.
**/
type GameInformation struct {
	MsgChan *GameMsgChan // message channel to communicate with game

	// --- metadata ---

	// players in this game instance
	players map[uuid.UUID]PlayerState
}
