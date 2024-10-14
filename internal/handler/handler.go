package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
)

type User interface {
	CreateNewUser(ctx context.Context, registerUser models.UserLogin) (models.JWTToken, error)
	LoginUser(ctx context.Context, loginUser models.UserLogin) (models.JWTToken, error)
}

type Order interface {
	Add(ctx context.Context, userID uint64, orderID string) error
	GetAll(ctx context.Context, userID uint64) ([]models.Order, error)
}

type Balance interface {
	GetBalance(ctx context.Context, userID uint64) (models.Balance, error)
	Withdraw(ctx context.Context, userID uint64, order string, sum float32) error
	GetWithdraws(ctx context.Context, userID uint64) ([]models.GetWithdraw, error)
}

type Handler struct {
	user    User
	order   Order
	balance Balance
	cfg     *config.Config
}

func New(
	user User,
	order Order,
	balance Balance,
	cfg *config.Config,
) *Handler {
	return &Handler{
		user:    user,
		order:   order,
		balance: balance,
		cfg:     cfg,
	}
}

func (h Handler) getUserID(req *http.Request) (uint64, error) {
	fmt.Println(req, "123123123")
	userIDStr := fmt.Sprintf("%v", req.Context().Value("user_id"))

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return userID, errors.New("user not found")
	}

	if userID <= 0 {
		return userID, errors.New("user not found")
	}

	return userID, nil
}
