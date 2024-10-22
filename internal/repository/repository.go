package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sixojke/test-astral/domain"
)

type User interface {
	Create(login, password string) error
	GetByCredentials(login, password string) (string, error)
	AddSession(session domain.Session) error
	GetUserIdBySession(session string) (userId string, err error)
	DeleteSession(session string) error
	GetUserIdByLogin(login string) (string, error)
}

type Document interface {
	Create(document *domain.Document, userId string) error
	GetCurrentUserDocuments(currentUserId string, params *domain.FilterParams) (*[]domain.Document, error)
	GetOtherUserDocuments(userId string, currentUserId string, params *domain.FilterParams) (*[]domain.Document, error)
	GetById(documentId, userId string) (*domain.Document, error)
	CheckById(documentId, userId string) (bool, error)
	Delete(documentId, userId string) (filePath string, err error)
}

type Deps struct {
	Postgres *sqlx.DB
}

type Repository struct {
	User
	Document
}

func NewService(deps *Deps) *Repository {
	return &Repository{
		NewUserPostgres(deps.Postgres),
		NewDocumentPostgres(deps.Postgres),
	}
}
