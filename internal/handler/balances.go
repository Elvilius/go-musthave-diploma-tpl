package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/utils"
)

func (h *Handler) GetBalance(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	ctx := req.Context()
	userID, err := h.getUserID(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusUnauthorized)
		return
	}

	balance, err := h.balance.GetBalance(ctx, userID)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	responseRaw, err := json.Marshal(balance)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Write(responseRaw)
}

func (h *Handler) Withdraw(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, err := h.getUserID(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	withdraw := models.Withdraw{}
	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok := utils.CheckLuhn(withdraw.Order); !ok {
		http.Error(resp, "invalid order ID", http.StatusUnprocessableEntity)
		return
	}

	err = h.balance.Withdraw(ctx, userID, withdraw.Order, withdraw.Sum)
	if err != nil {
		if errors.Is(err, models.ErrInsufficientBalance) {
			http.Error(resp, "insufficient balance", http.StatusPaymentRequired)
			return
		}
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWithdraw(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	ctx := req.Context()
	userID, err := h.getUserID(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusUnauthorized)
		return
	}

	withdraws, err := h.balance.GetWithdraws(ctx, userID)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(withdraws)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Write(responseData)
}
