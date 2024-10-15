package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func New(ctx context.Context, dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 4; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.PingContext(ctx)
			if err == nil {
				return db, nil
			}

			if isRetriableError(err) {
				fmt.Printf("Attempt %d: retriable error encountered, retrying...\n", i+1)
			} else {
				return nil, fmt.Errorf("non-retriable error: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to open connection: %w", err)
		}

		waitTime := time.Duration(1+2*i) * time.Second
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
}

func isRetriableError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case pgerrcode.ConnectionException, pgerrcode.ConnectionDoesNotExist,
			pgerrcode.ConnectionFailure, pgerrcode.SQLClientUnableToEstablishSQLConnection,
			pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
			pgerrcode.TransactionResolutionUnknown:
			return true
		}
	}
	return false
}
