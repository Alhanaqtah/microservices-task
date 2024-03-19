package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"user-managment-service/internal/config"
	"user-managment-service/internal/models"
	"user-managment-service/internal/storage/storage"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	pool *pgx.Conn
}

func New(cfg config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgx.Connect(context.Background(), cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	/* _, err = pool.Exec(context.Background(), `
	CREATE TYPE ROLE AS ENUM('user', 'admin', 'moderator');

	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(100) DEFAULT '',
		surname VARCHAR(100) DEFAULT '',
		username TEXT UNIQUE NOT NULL,
		pass_hash BYTEA NOT NULL,
		phone_number VARCHAR(20) DEFAULT '',
		email VARCHAR(255) DEFAULT '',
		role ROLE DEFAULT 'user',
		group_id INTEGER REFERENCES groups(id),
		image_s3_path TEXT DEFAULT '',
		is_blocked BOOLEAN DEFAULT false,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		modified_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	} */

	return &Storage{pool: pool}, nil
}

func (s *Storage) UserByUUID(ctx context.Context, username string) (*models.User, error) {
	const op = "storage.postgres.UserByUUID"

	row := s.pool.QueryRow(ctx, `
		SELECT
			id,
			name,
			surname,
			username,
			pass_hash,
			phone_number,
			email,
			role,
			group_id,
			image_s3_path,
			is_blocked,
			created_at,
			modified_at
		FROM users WHERE id=$1`, username,
	)

	var groupID sql.NullInt64
	var user models.User
	err := row.Scan(
		&user.UUID,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.PassHash,
		&user.PhoneNumber,
		&user.Email,
		&user.Role,
		&groupID,
		&user.ImageS3Path,
		&user.IsBlocked,
		&user.CreatedAt,
		&user.ModifiedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if groupID.Valid {
		user.GroupID = groupID.Int64
	}

	return &user, nil
}

func (s *Storage) UserByName(ctx context.Context, username string) (*models.User, error) {
	const op = "storage.postgres.UserByName"

	row := s.pool.QueryRow(ctx, `
		SELECT
			id,
			name,
			surname,
			username,
			pass_hash,
			phone_number,
			email,
			role,
			group_id,
			image_s3_path,
			is_blocked,
			created_at,
			modified_at
		FROM users WHERE username=$1`, username,
	)

	var groupID sql.NullInt64
	var user models.User
	err := row.Scan(
		&user.UUID,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.PassHash,
		&user.PhoneNumber,
		&user.Email,
		&user.Role,
		&groupID,
		&user.ImageS3Path,
		&user.IsBlocked,
		&user.CreatedAt,
		&user.ModifiedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if groupID.Valid {
		user.GroupID = groupID.Int64
	}

	return &user, nil
}

func (s *Storage) CreateNewUser(ctx context.Context, username string, passHash []byte) (string, error) {
	const op = "storage.postgres.CreateNewUser"

	var id string

	err := s.pool.QueryRow(ctx, `INSERT INTO users (username, pass_hash) VALUES ($1, $2) RETURNING id`, username, passHash).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
