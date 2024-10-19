package service

import (
	"github.com/sixojke/test-astral/internal/config"
	"github.com/sixojke/test-astral/internal/repository"
	"github.com/sixojke/test-astral/pkg/hash"
)

type User interface {
	SignUp(adminToken, login, password string) error
}

type Deps struct {
	Repository *repository.Repository
	Config     *config.Config
	Hasher     hash.PasswordHasher
}

type Service struct {
	User
}

func NewService(deps *Deps) *Service {
	return &Service{
		NewUserService(deps.Repository.User, deps.Hasher, deps.Config.Authorization.AdminToken),
	}
}
