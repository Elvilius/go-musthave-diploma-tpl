package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/lib/pq"
)

func (s *Store) AddNewOrder(ctx context.Context, userID uint64, orderID string) (models.Order, error) {
	var order models.Order

	tx, err := s.DB.Begin()
	if err != nil {
		return order, err
	}

	if order, err := s.GetOrder(ctx, orderID); err == nil {
		if order.UserID != int(userID) {
			return order, models.ErrOrderAlreadyUploadedByAnotherUser
		}
		return order, models.ErrOrderInProcessed
	}

	args := []interface{}{orderID, userID, models.NEW, time.Now().Format(time.RFC3339)}
	query := "INSERT INTO orders (number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)  RETURNING number, user_id, status, uploaded_at"
	result := tx.QueryRowContext(ctx, query, args...)

	errExec := result.Scan(&order.Number, &order.UserID, &order.Status, &order.UploadedAt)
	if errExec != nil {
		if pqErr, ok := errExec.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return order, models.ErrOrderExist
			}
		}
		tx.Rollback()
		return order, errExec
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return order, err
	}

	return order, nil
}

func (s *Store) GetAllOrders(ctx context.Context, userID uint64) ([]models.Order, error) {
	var orders []models.Order
	result, err := s.DB.QueryContext(ctx, "SELECT number, status, accrual, uploaded_at, user_id FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var order models.Order
		if err := result.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Store) GetPendingOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	tx, err := s.DB.Begin()
	if err != nil {
		return orders, err
	}

	statuses := []string{"NEW", "REGISTERED", "PROCESSING"}

	placeholders := make([]string, len(statuses))
	args := make([]interface{}, len(statuses))
	for i, status := range statuses {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = status
	}

	query := fmt.Sprintf("SELECT number, status, accrual, uploaded_at FROM orders WHERE status IN (%s)", strings.Join(placeholders, ", "))

	result, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.Err() != nil {
		tx.Rollback()
		return nil, err
	}

	for result.Next() {
		var order models.Order
		err := result.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Store) UpdateOrder(ctx context.Context, order models.Order) error {
	var userID uint64

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	args := []interface{}{order.Status, order.Accrual, order.Number}
	query := "UPDATE orders set status = $1, accrual = $2 WHERE number = $3 RETURNING user_id"
	result := tx.QueryRowContext(ctx, query, args...)

	err = result.Scan(&userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.updateBalance(ctx, tx, userID, order.Accrual)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetOrder(ctx context.Context, number string) (models.Order, error) {
	args := []interface{}{number}
	query := "SELECT number, status, accrual, uploaded_at, user_id FROM orders WHERE number = $1"
	result := s.DB.QueryRowContext(ctx, query, args...)

	var order models.Order

	err := result.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)

	return order, err
}
