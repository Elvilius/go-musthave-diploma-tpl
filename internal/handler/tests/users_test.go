package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/balances"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/handler"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/orders"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/users"
	mocks_users "github.com/Elvilius/go-musthave-diploma-tpl.git/internal/users/mocks"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_RegisterUser(t *testing.T) {
	type want struct {
		status int
		token  string
	}
	tests := []struct {
		name     string
		body     models.UserLogin
		want     want
		mockFunk func(userStore *mocks_users.MockStorer)
	}{
		{
			name: "success create new user",
			body: models.UserLogin{Login: "test", Password: "123456"},
			mockFunk: func(userStore *mocks_users.MockStorer) {
				userStore.EXPECT().CreateUser(gomock.Any(), "test", gomock.Any()).Return(1, nil)
			},
			want: want{
				status: http.StatusOK,
				token:  "secret",
			},
		},
		{
			name: "validate error",
			body: models.UserLogin{},
			want: want{
				status: http.StatusBadRequest,
				token:  "",
			},
		},
		{
			name: "error user exists",
			body: models.UserLogin{Login: "test", Password: "123"},
			mockFunk: func(userStore *mocks_users.MockStorer) {
				userStore.EXPECT().CreateUser(gomock.Any(), "test", gomock.Any()).Return(1, models.ErrUserExists)
			},
			want: want{
				status: http.StatusConflict,
				token:  "",
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userStore := mocks_users.NewMockStorer(ctrl)
		token := jwt.NewMockToken()
		userService := users.New(userStore, token, nil)

		orderService := orders.New(nil, nil, nil, nil)
		balanceService := balances.New(nil)

		h := handler.New(userService, orderService, balanceService, nil)

		router := chi.NewRouter()
		router.Post("/api/user/register", h.RegisterUser)

		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(jsonData))

			if tt.mockFunk != nil {
				tt.mockFunk(userStore)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			assert.Equal(t, tt.want.token, result.Header.Get("Authorization"))

			result.Body.Close()
		})
	}
}

func TestHandler_Login(t *testing.T) {
	type want struct {
		status int
		token  string
	}
	tests := []struct {
		name     string
		body     models.UserLogin
		want     want
		mockFunk func(userStore *mocks_users.MockStorer)
	}{
		{
			name: "success login",
			body: models.UserLogin{Login: "test", Password: "test"},
			mockFunk: func(userStore *mocks_users.MockStorer) {
				userStore.EXPECT().GetUserByLogin(gomock.Any(), "test").Return(models.User{Login: "test", PasswordHash: "$2a$10$EmS6L/RlhGUisab3AAVAUuNShZCmr838QqbekeJEMv56MYbAgkDoC"}, nil)
			},
			want: want{
				status: http.StatusOK,
				token:  "secret",
			},
		},
		{
			name: "validate error",
			body: models.UserLogin{},
			want: want{
				status: http.StatusBadRequest,
				token:  "",
			},
		},
		{
			name: "error not valid password",
			body: models.UserLogin{Login: "test", Password: "123"},
			mockFunk: func(userStore *mocks_users.MockStorer) {
				userStore.EXPECT().GetUserByLogin(gomock.Any(), "test").Return(models.User{Login: "test", PasswordHash: "$2a$10$EmS6L/RlhGUisab3AAVAUuNShZCmr838QqbekeJEMv56MYbAgkDoC"}, nil)
			},
			want: want{
				status: http.StatusUnauthorized,
				token:  "",
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userStore := mocks_users.NewMockStorer(ctrl)
		token := jwt.NewMockToken()
		userService := users.New(userStore, token, nil)

		orderService := orders.New(nil, nil, nil, nil)
		balanceService := balances.New(nil)

		h := handler.New(userService, orderService, balanceService, nil)

		router := chi.NewRouter()
		router.Post("/api/user/login", h.LoginUser)

		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(jsonData))

			if tt.mockFunk != nil {
				tt.mockFunk(userStore)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			assert.Equal(t, tt.want.token, result.Header.Get("Authorization"))

			result.Body.Close()
		})
	}
}
