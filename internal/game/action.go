package game

/**
* Each round in the game, for each player, is an "Event".
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

func (e *GameEvent) initiateEvent() {

}
