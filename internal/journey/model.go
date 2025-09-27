package journey

import (
	"time"

	"github.com/google/uuid"
)

type Journey struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type PlayerJourney struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	PlayerProfileID uuid.UUID  `db:"player_profile_id" json:"player_profile_id"`
	JourneyID       uuid.UUID  `db:"journey_id" json:"journey_id"`
	JoinedAt        time.Time  `db:"joined_at" json:"joined_at"`
	CompletedAt     *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	Status          string     `db:"status" json:"status"`
}