package user

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo Repository
}

type Repository interface {
	Create(request SignUpRequest) (*User, error)
	GetByUsername(username string) (*User, error)
	GetById(userId uuid.UUID) (*User, error)
	GetProfile(userId uuid.UUID) (*PlayerProfile, error)
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) SignUp(request SignUpRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	request.Password = string(hashedPassword)

	_, err = s.repo.Create(request)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *service) SignIn(request SignInRequest) (*AuthResponse, error) {
	user, err := s.repo.GetByUsername(request.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *service) RefreshToken(request RefreshTokenRequest) (*AuthResponse, error) {
	claims, err := ValidateToken(request.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	userId, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token")
	}

	user, err := s.repo.GetById(userId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	accessToken, err := GenerateAccessToken(userId.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateRefreshToken(userId.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *service) GetProfile(userId uuid.UUID) (*PlayerProfile, error) {
	profile, err := s.repo.GetProfile(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile, nil
}
