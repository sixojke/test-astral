package migrations

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/sixojke/test-astral/internal/config"
	"github.com/sixojke/test-astral/pkg/utils"
)

func MigratePostgres(cfg config.Postgres) error {
	err := migrateUp(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode), "/postgres")
	if err != nil {
		return err
	}

	return nil
}

func migrateUp(conn, dir string) error {
	path, err := utils.CustomPath("/schema")
	if err != nil {
		return err
	}
	m, err := migrate.New(
		"file:"+path+dir,
		conn,
	)
	if err != nil {
		return fmt.Errorf("create migration: %s", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("load migration: %s", err)
		}
	}

	return nil
}
