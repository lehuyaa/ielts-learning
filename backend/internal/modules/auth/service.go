package auth

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"ielts-learning/backend/internal/models"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

var (
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)

type Service struct {
	repository Repository
	jwtManager sharedjwt.Manager
}

func NewService(repository Repository, jwtManager sharedjwt.Manager) Service {
	return Service{
		repository: repository,
		jwtManager: jwtManager,
	}
}

func (s Service) Register(req RegisterRequest) (AuthResponse, error) {
	req = NormalizeRegisterRequest(req)
	email := req.Email
	exists, err := s.repository.EmailExists(email)
	if err != nil {
		return AuthResponse{}, err
	}
	if exists {
		return AuthResponse{}, ErrEmailAlreadyUsed
	}

	username := req.Username
	if username != nil {
		exists, err := s.repository.UsernameExists(*username)
		if err != nil {
			return AuthResponse{}, err
		}
		if exists {
			return AuthResponse{}, ErrUsernameAlreadyUsed
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("hash password: %w", err)
	}

	targetBand := 7.0
	if req.TargetBand != nil {
		targetBand = *req.TargetBand
	}

	user := models.User{
		Email:        email,
		Name:         req.Name,
		Username:     username,
		PasswordHash: string(passwordHash),
		Role:         models.UserRoleUser,
		TargetBand:   targetBand,
		Timezone:     "UTC",
		Locale:       "en",
	}

	if err := s.repository.Create(&user); err != nil {
		return AuthResponse{}, err
	}

	return s.authResponse(user)
}

func (s Service) Login(req LoginRequest) (AuthResponse, error) {
	req = NormalizeLoginRequest(req)
	user, err := s.repository.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return AuthResponse{}, ErrInvalidCredentials
		}
		return AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	return s.authResponse(user)
}

func (s Service) Me(userID uint) (UserResponse, error) {
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return UserResponse{}, err
	}

	return toUserResponse(user), nil
}

func (s Service) authResponse(user models.User) (AuthResponse, error) {
	token, err := s.jwtManager.Generate(user.ID, user.Email, string(user.Role))
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User:        toUserResponse(user),
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := strings.ToLower(strings.TrimSpace(*value))
	if normalized == "" {
		return nil
	}

	return &normalized
}
