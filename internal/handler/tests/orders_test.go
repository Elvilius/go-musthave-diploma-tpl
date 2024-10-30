package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/balances"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/handler"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/orders"
	mocks_orders "github.com/Elvilius/go-musthave-diploma-tpl.git/internal/orders/mocks"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/users"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_AddNewOrder(t *testing.T) {
	type want struct {
		status int
	}
	tests := []struct {
		name     string
		body     string
		want     want
		callMock bool
	}{
		{
			name: "success create new order",
			body: "48113687317151",

			want: want{
				status: http.StatusAccepted,
			},
			callMock: true,
		},

		{
			name: "error create new order",
			body: "332332",

			want: want{
				status: http.StatusUnprocessableEntity,
			},
			callMock: false,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		orderStore := mocks_orders.NewMockStorer(ctrl)
		token := jwt.NewMockToken()
		userService := users.New(nil, token, nil)

		orderService := orders.New(orderStore, nil, nil, nil)
		balanceService := balances.New(nil)

		h := handler.New(userService, orderService, balanceService, nil)

		if tt.callMock {
			orderStore.EXPECT().AddNewOrder(gomock.Any(), gomock.Any(), tt.body).Return(models.Order{Number: tt.body}, nil)
		}
		router := chi.NewRouter()
		router.Post("/api/user/orders", h.AddNewOrder)

		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewReader([]byte(tt.body)))
			request = SetAuth(request)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)

			result.Body.Close()
		})
	}
}

