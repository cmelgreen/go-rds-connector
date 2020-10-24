package main

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"
)

const (
	timeout = 5
	sqlDriver = "pgx"
)



// Database abstracts sqlx connection
type Database struct {
	*sqlx.DB
}

// ConnectToDB creates a db connection with a predefined timeout
func ConnectToDB(ctx context.Context) (*Database, error) {
	ctx, cancelFn := context.WithTimeout(ctx, timeout*time.Second)
	defer cancelFn()

	config, err := ConfigString(ctx)
	if err != nil {
		return &Database{}, err
	}

	conn, err := sqlx.ConnectContext(ctx, sqlDriver, config)
	if err != nil {
		return &Database{}, err
	}

	return &Database{conn}, nil
}

// Connected pings server and returns bool response status
func (db *Database) Connected(ctx context.Context) bool {
	if db == nil {
		return false
	}

	err := db.PingContext(ctx)

	if err != nil {
		return false
	}

	return true
}