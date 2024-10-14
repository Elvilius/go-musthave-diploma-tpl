package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
)

func (h *Handler) RegisterUser(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var registerUser models.UserLogin

	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&registerUser)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	var validateErrors []string
	if registerUser.Login == "" {
		validateErrors = append(validateErrors, "login must required")
	}

	if len(registerUser.Password) <= 0 {
		validateErrors = append(validateErrors, "password must required")
	}

	if len(validateErrors) != 0 {
		http.Error(resp, strings.Join(validateErrors, "\n"), http.StatusBadRequest)
		return
	}

	authToken, err := h.user.CreateNewUser(ctx, registerUser)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			http.Error(resp, "user exists with this login", http.StatusConflict)
			return
		}
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Authorization", string(authToken))
	resp.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginUser(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var registerUser models.UserLogin

	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&registerUser)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if registerUser.Login == "" {
		http.Error(resp, "login must required", http.StatusBadRequest)
		return
	}
	authToken, err := h.user.LoginUser(ctx, registerUser)
	if err != nil {
		if errors.Is(err, models.ErrUserPasswordNotValid) {
			http.Error(resp, "", http.StatusUnauthorized)
			return
		}
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Authorization", string(authToken))
	resp.WriteHeader(http.StatusOK)
}
