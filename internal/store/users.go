package store

import (
	"context"
	"database/sql"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/lib/pq"
)

func (s *Store) CreateUser(ctx context.Context, login string, passwordHash string) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}

	userID, err := s.createUser(ctx, tx, login, passwordHash)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, models.ErrUserExists
			}
		}
		return 0, err
	}

	err = s.CreateBalance(ctx, tx, userID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return userID, nil
}

func (s *Store) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User

	res := s.DB.QueryRowContext(ctx, "SELECT id, login, password_hash from users WHERE login = $1", login)

	err := res.Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *Store) createUser(ctx context.Context, tx *sql.Tx, login, passwordHash string) (int, error) {
	var userID int

	args := []interface{}{login, passwordHash}
	query := "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id"
	result := tx.QueryRowContext(ctx, query, args...)

	if result.Err() != nil {
		return userID, result.Err()
	}

	err := result.Scan(&userID)

	return userID, err
}
