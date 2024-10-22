package service

import (
	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/internal/config"
	"github.com/sixojke/test-astral/internal/repository"
	"github.com/sixojke/test-astral/pkg/auth"
	"github.com/sixojke/test-astral/pkg/hash"
)

type User interface {
	SignUp(adminToken, login, password string) error
	SignIn(login, password string) (accessToken string, err error)
	GetUserIdByToken(token string) (userId string, err error)
	DeleteSession(token string) error
}

type Document interface {
	Create(document *domain.Document, userId string) error
	GetByUser(userLogin, currentUserId string, params *domain.FilterParams) (*[]domain.Document, error)
	GetById(documentId, userId string) (*domain.Document, error)
	CheckById(documentId, userId string) (bool, error)
	Delete(documentId, userId string) error
}

type Deps struct {
	Repository   *repository.Repository
	Config       *config.Config
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
}

type Service struct {
	User
	Document
}

func NewService(deps *Deps) *Service {
	return &Service{
		NewUserService(deps.Repository.User, deps.Hasher, deps.Config.Authorization, deps.TokenManager),
		NewDocumentService(deps.Repository.Document, deps.Repository.User),
	}
}
