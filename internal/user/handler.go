package user

import "github.com/google/uuid"

type Handler struct {
	service Service
}

type Service interface {
	SignUp(request SignUpRequest) error
	SignIn(request SignInRequest) (*User, error)
	GetProfile(userId uuid.UUID) (*PlayerProfile, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
