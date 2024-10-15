package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
)

func (s *Store) CreateBalance(ctx context.Context, tx *sql.Tx, userID int) error {
	args := []interface{}{userID}
	query := "INSERT INTO balances (user_id) VALUES ($1)"
	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) GetBalance(ctx context.Context, userID uint64) (models.Balance, error) {
	var balance models.Balance
	tx, err := s.DB.Begin()
	if err != nil {
		return balance, err
	}

	args := []interface{}{userID}
	result := tx.QueryRowContext(ctx, "SELECT current_balance, withdrawn from balances WHERE user_id = $1", args...)

	err = result.Scan(&balance.CurrentBalance, &balance.Withdrawn)
	if err != nil {
		err := tx.Rollback()
		return balance, err
	}

	if err := tx.Commit(); err != nil {
		return balance, err
	}

	fmt.Println(float32(balance.CurrentBalance))
	return balance, nil
}

func (s *Store) checkBalance(ctx context.Context, tx *sql.Tx, userID uint64, balance float32) error {
	args := []interface{}{userID}
	result := tx.QueryRowContext(ctx, "SELECT current_balance from balances WHERE user_id = $1", args...)

	var currentBalance float32
	err := result.Scan(&currentBalance)

	if currentBalance-balance < 0 {
		return models.ErrInsufficientBalance
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) withdraw(ctx context.Context, tx *sql.Tx, userID uint64, sum float32, order string) error {
	args := []interface{}{userID, order, sum}
	query := "INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)"
	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) updateBalance(ctx context.Context, tx *sql.Tx, userID uint64, sum float32) error {
	var query string
	args := []interface{}{sum, userID}

	if sum > 0 {
		query = "UPDATE balances SET current_balance = current_balance + $1 WHERE user_id = $2"
	} else if sum < 0 {
		query = `
            UPDATE balances 
            SET 
                current_balance = current_balance - $1,
                withdrawn = withdrawn + $1 
            WHERE user_id = $2`
		args[0] = -sum
	}

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) Withdraw(ctx context.Context, userID uint64, order string, sum float32) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	err = s.checkBalance(ctx, tx, userID, sum)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.updateBalance(ctx, tx, userID, -sum)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.withdraw(ctx, tx, userID, sum, order)
	if err != nil {

		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetWithdraws(ctx context.Context, userID uint64) ([]models.GetWithdraw, error) {
	var withdraws []models.GetWithdraw
	tx, err := s.DB.Begin()
	if err != nil {
		return withdraws, err
	}

	args := []interface{}{userID}
	result, err := tx.QueryContext(ctx, "SELECT order_number, sum, processed_at from withdrawals WHERE user_id = $1 ORDER BY processed_at DESC", args...)
	if err != nil {
		tx.Rollback()
		return withdraws, err
	}
	defer result.Close()

	for result.Next() {
		var withdraw models.GetWithdraw
		err := result.Scan(&withdraw.OrderNumber, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			tx.Rollback()
			return withdraws, err
		}
		withdraws = append(withdraws, withdraw)
	}

	if err = result.Err(); err != nil {
		tx.Rollback()
		return withdraws, err
	}

	if err = tx.Commit(); err != nil {
		return withdraws, err
	}

	return withdraws, nil
}
