package game

import (
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

/**
* Factory for creating games with default settings and player state
* conversion.
**/
type GameFactory struct {
	defaultRound       int
	defaultPlayerState DefaultPlayerState
}

type InitialConditions struct {
	RoundDefault   int
	PlayerDefaults struct {
		DefaultGold  int
		DefaultItems []model.Item
	}
}

type DefaultPlayerState struct {
	defaultGold       int
	defaultItems      []model.Item
	defaultAttributes model.Attributes
}

func NewGameFactory(initialConditions InitialConditions) *GameFactory {

	return &GameFactory{
		defaultRound: initialConditions.RoundDefault,
		defaultPlayerState: DefaultPlayerState{
			defaultGold:  initialConditions.PlayerDefaults.DefaultGold,
			defaultItems: initialConditions.PlayerDefaults.DefaultItems,
			defaultAttributes: model.Attributes{
				Strength:     1,
				Intelligence: 1,
				Wisdom:       1,
				Agility:      1,
				Vitality:     1,
				Faith:        1,
				Charisma:     1,
			},
		},
	}
}

/**
* Core game creation, based on factory settings and converting all passed
* in players into "PlayerState"s.
**/
func (f *GameFactory) CreateGame(players []*model.Player) *model.Game {
	// convert Players to PlayerState
	playerStates := make(map[uuid.UUID]*model.PlayerState)

	for _, player := range players {
		playerStates[player.ID] = &model.PlayerState{
			Player:    *player,
			Gold:      f.defaultPlayerState.defaultGold,
			Inventory: f.defaultPlayerState.defaultItems,
		}
	}

	return &model.Game{
		ID:      uuid.New(),
		MsgChan: make(chan string),
		Round:   f.defaultRound,
		Players: playerStates,
	}
}
