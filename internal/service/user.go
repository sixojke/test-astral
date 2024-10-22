package service

import (
	"errors"
	"time"

	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/internal/config"
	"github.com/sixojke/test-astral/internal/repository"
	"github.com/sixojke/test-astral/pkg/auth"
	"github.com/sixojke/test-astral/pkg/hash"
	"github.com/sixojke/test-astral/pkg/logger"
)

type UserService struct {
	repo         repository.User
	hasher       hash.PasswordHasher
	authConfig   config.Authorization
	tokenManager auth.TokenManager
}

func NewUserService(repo repository.User, hasher hash.PasswordHasher, authConfig config.Authorization,
	tokenManager auth.TokenManager) *UserService {
	return &UserService{
		repo:         repo,
		hasher:       hasher,
		authConfig:   authConfig,
		tokenManager: tokenManager,
	}
}

func (s *UserService) SignUp(adminToken, login, password string) error {
	if s.authConfig.AdminToken != adminToken {
		return domain.ErrInvalidToken
	}

	pswdHash, err := s.hasher.Hash(password)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		return err
	}

	if err := s.repo.Create(login, pswdHash); err != nil {
		if !errors.Is(err, domain.ErrLoginIsBusy) {
			logger.Errorf("failed to sign up user: %v", err)
		}

		return err
	}

	return nil
}

func (s *UserService) SignIn(login, password string) (accessToken string, err error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}

	userId, err := s.repo.GetByCredentials(login, passwordHash)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			logger.Errorf("failed to get user by credentials: %v", err)
		}

		return "", err
	}

	return s.createSession(userId)
}

func (s *UserService) createSession(userId string) (accessToken string, err error) {
	accessToken, err = s.tokenManager.NewJWT(userId, s.authConfig.JWT.AccessTokenTTL)
	if err != nil {
		logger.Errorf("failed to create access token: %v", err)
		return accessToken, err
	}

	if err = s.repo.AddSession(domain.Session{
		UserId:      userId,
		AccessToken: accessToken,
		ExpiresAt:   time.Now().Add(s.authConfig.JWT.AccessTokenTTL),
	}); err != nil {
		logger.Errorf("failed to add session: %v", err)
		return "", err
	}

	return accessToken, err
}

func (s *UserService) GetUserIdByToken(token string) (userId string, err error) {
	return s.repo.GetUserIdBySession(token)
}

func (s *UserService) DeleteSession(token string) error {
	return s.repo.DeleteSession(token)
}
