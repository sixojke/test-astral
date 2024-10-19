package service

import (
	"errors"

	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/internal/repository"
	"github.com/sixojke/test-astral/pkg/hash"
	"github.com/sixojke/test-astral/pkg/logger"
)

type UserService struct {
	repo       repository.User
	hasher     hash.PasswordHasher
	adminToken string
}

func NewUserService(repo repository.User, hasher hash.PasswordHasher, adminToken string) *UserService {
	return &UserService{
		repo:       repo,
		hasher:     hasher,
		adminToken: adminToken,
	}
}

func (s *UserService) SignUp(adminToken, login, password string) error {
	if s.adminToken != adminToken {
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
