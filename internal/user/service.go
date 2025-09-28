package user

import "github.com/google/uuid"

type service struct {
	repo Repository
}

type Repository interface {
	Create(request SignUpRequest)
	GetProfile(userId uuid.UUID) (*PlayerProfile, error)
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}
