package model

import "github.com/google/uuid"

type GameMsgChan = chan string

type Time struct {
	day  int
	hour int
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
	strength     int
	intelligence int
	wisdom       int
	agility      int
	vitality     int
	faith        int
	charisma     int
}

type Item struct {
	id   uuid.UUID
	name string
}

type Skill struct {
	id   uuid.UUID
	name string
}

type PlayerState struct {
	Player
	// the stage the game match has reached for any player
	time map[uuid.UUID]Time

	// hero stuff
	heroInfo Hero
}

type Hero struct {
	class      HeroClass  // mage
	attributes Attributes //
	skills     []Skill
	items      []Item
}

type Follower struct {
	heroInfo Hero
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
