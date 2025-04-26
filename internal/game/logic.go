package game

import (
	"fmt"
	"github.com/google/uuid"
)

func (g *Game) RemovePlayer(id uuid.UUID) error {
	_, exists := g.Players[id]

	if exists {
		// remove player
		delete(g.Players, id)
		return nil
	}

	return fmt.Errorf("Player with id %s does not exist in this game:", id)
}
