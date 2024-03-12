package postgres

import (
	"database/sql"
	"fmt"
	"user-managment-service/internal/config"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.NetLoc, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TYPE ROLE AS ENUM ('user', 'admin', 'moderator');
	
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name TEXT NOT NULL,
			surname TEXT NOT NULL,
			username TEXT UNIQUE NOT NULL,
			pass_hash BLOB NOT NULL,
			phone_number TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			role ROLE NOT NULL,
			group INTEGER REFERENCES groups(id),
			image_s3_path TEXT NOT NULL,
			is_blocked BOOL NOT NULL,
			created_at DATETIME NOT NULL,
			modified_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS groups (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME NOT NULL
		);
	`)

	return &Storage{db: db}, nil
}
