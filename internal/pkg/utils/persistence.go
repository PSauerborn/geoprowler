package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PostgresPersistence struct {
	// define connection settions for persistence
	DatabaseURL string
	// define postgres session/connection pool
	Session *pgxpool.Pool
}

// function used to close persistence
// connections. should be called in defer
// block/statement
func (db *PostgresPersistence) Close() {
	db.Session.Close()
}

// function used to connect/generate connection
// pool to postgres instance
func (db *PostgresPersistence) Connect() error {
	log.Debug("generating postgres connection pool...")
	// connect to postgres server and set session in persistence
	conn, err := pgxpool.Connect(context.Background(), db.DatabaseURL)
	if err != nil {
		log.Error(fmt.Errorf("error connecting to postgres service: %+v", err))
		return err
	}
	db.Session = conn
	return err
}
