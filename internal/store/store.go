package store

import (
	"database/sql"
)

const DuplicateCode = "23505"

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}
