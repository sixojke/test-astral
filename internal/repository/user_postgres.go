package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/pkg/logger"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (r *UserPostgres) Create(login, password string) error {
	query := `
		INSERT INTO users (
			login,
			password_hash
		) 
		VALUES
			($1, $2)
	`

	if _, err := r.db.Exec(query, login, password); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrLoginIsBusy
		}

		logger.Errorf("failed to create user: %v", err)
		return err
	}

	return nil
}
