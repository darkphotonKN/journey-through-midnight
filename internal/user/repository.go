package user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(request SignUpRequest) (*User, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var user User
	userQuery := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, email_verified, is_active, created_at, updated_at
	`
	err = tx.QueryRowx(userQuery, request.Email, request.Password).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	profileQuery := `
		INSERT INTO player_profiles (user_id, username)
		VALUES ($1, $2)
		RETURNING id, user_id, username, reputation, total_playtime, created_at, updated_at
	`
	var profile PlayerProfile
	err = tx.QueryRowx(profileQuery, user.ID, request.Username).StructScan(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create player profile: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

func (r *repository) GetByUsername(username string) (*User, error) {
	var user User
	query := `
		SELECT u.id, u.email, u.password_hash, u.email_verified, u.is_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN player_profiles pp ON u.id = pp.user_id
		WHERE pp.username = $1
	`
	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

func (r *repository) GetById(userId uuid.UUID) (*User, error) {
	var user User
	query := `
		SELECT id, email, password_hash, email_verified, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.Get(&user, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *repository) GetProfile(userId uuid.UUID) (*PlayerProfile, error) {
	var profile PlayerProfile
	query := `
		SELECT id, user_id, username, reputation, total_playtime, created_at, updated_at
		FROM player_profiles
		WHERE user_id = $1
	`
	err := r.db.Get(&profile, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	return &profile, nil
}
