package repository

import "github.com/jmoiron/sqlx"

type User interface {
	Create(login, password string) error
}

type Deps struct {
	Postgres *sqlx.DB
}

type Repository struct {
	User
}

func NewService(deps *Deps) *Repository {
	return &Repository{
		NewUserPostgres(deps.Postgres),
	}
}
