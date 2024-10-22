package repository

import (
	"database/sql"
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
	logger.Debugf("create user: params=[login=%v, password=%v]", login, password)

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

func (r *UserPostgres) GetByCredentials(login, password string) (string, error) {
	logger.Debugf("get user by credentials: params=[login=%v password=%v]", login, password)

	query := `
		SELECT 
			id
		FROM users 
		WHERE login = $1 AND password_hash = $2`

	var id string
	if err := r.db.Get(&id, query, login, password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}

		logger.Errorf("failed to get user by credentials: %v", err)
		return "", err
	}

	return id, nil
}

func (r *UserPostgres) AddSession(session domain.Session) error {
	logger.Debugf("add session: params[session=%v]", session)

	query := `
		INSERT INTO tokens (
			user_id,
			token,
			expires_at
		) VALUES
			($1, $2, $3)
	`

	if _, err := r.db.Exec(query, session.UserId, session.AccessToken, session.ExpiresAt); err != nil {
		logger.Errorf("failed to insert token: %v", err)
		return err
	}

	return nil
}

func (r *UserPostgres) GetUserIdBySession(session string) (userId string, err error) {
	logger.Debugf("get userId by session: params[session=%v]", session)

	query := `
		SELECT 
			user_id
		FROM tokens
		WHERE 
			token = $1
			AND expires_at > now()
	`

	if err = r.db.Get(&userId, query, session); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}

		return "", err
	}

	return userId, nil
}

func (r *UserPostgres) DeleteSession(session string) error {
	logger.Debugf("delete session: params[session=%v]", session)
	query := `
	  DELETE FROM tokens
	  WHERE token = $1
	`

	if _, err := r.db.Exec(query, session); err != nil {
		logger.Errorf("failed to delete session: %v", err)
		return err
	}

	return nil
}

func (r *UserPostgres) GetUserIdByLogin(login string) (string, error) {
	logger.Debugf("get userId by login: params[login=%v]", login)

	query := `
		SELECT id
		FROM users
		WHERE login = $1
	`

	var userId string
	if err := r.db.Get(&userId, query, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}

		logger.Errorf("failed to get user_id by login: %v", err)
		return "", err
	}

	return userId, nil
}
