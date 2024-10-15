package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/balances"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/handler"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/mocks"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/orders"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

type contextKey string

const UserIDKey = contextKey("user_id")

func TestHandler_GetBalance(t *testing.T) {
	cfg := config.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userStore := mocks.NewMockUserStore(ctrl)
	orderStore := mocks.NewMockOrderStore(ctrl)
	balancesStore := mocks.NewMockBalancesStore(ctrl)
	token := mocks.NewMockToken()
	userService := users.New(userStore, token, cfg)
	orderService := orders.New(orderStore, nil, nil, nil)
	balanceService := balances.New(balancesStore)

	h := handler.New(userService, orderService, balanceService, cfg)

	balancesStore.EXPECT().GetBalance(gomock.Any(), uint64(1)).Return(models.Balance{CurrentBalance: 100, Withdrawn: 100}, nil)

	router := chi.NewRouter()
	router.Post("/api/user/balance", h.GetBalance)

	request := httptest.NewRequest(http.MethodPost, "/api/user/balance", nil)
	request = request.WithContext(context.WithValue(request.Context(), UserIDKey, 1))
	request.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjJ9.6Oz5eGuwTSWswdvgsxbhvDIBkd9YKzxJSyd9mg4auBM")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	result := w.Result()
	result.Body.Close()

	var balance models.Balance

	b, _ := io.ReadAll(result.Body)

	json.Unmarshal(b, &balance)

	assert.Equal(t, balance.CurrentBalance, float32(100))
	assert.Equal(t, balance.Withdrawn, float32(100))
	assert.Equal(t, result.StatusCode, http.StatusOK)
}
