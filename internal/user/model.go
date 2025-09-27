package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Email         string    `db:"email" json:"email"`
	PasswordHash  string    `db:"password_hash" json:"-"`
	EmailVerified bool      `db:"email_verified" json:"email_verified"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type PlayerProfile struct {
	ID            uuid.UUID `db:"id" json:"id"`
	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	Username      string    `db:"username" json:"username"`
	Reputation    int       `db:"reputation" json:"reputation"`
	TotalPlaytime int       `db:"total_playtime" json:"total_playtime"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type SignUpRequest struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Email    string `db:"Email" json:"Email"`
}

type SignInRequest struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}
