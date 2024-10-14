package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/utils"
)

func (h *Handler) AddNewOrder(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, err := h.getUserID(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusUnauthorized)
		return
	}

	inputData, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	orderID := string(inputData)

	if ok := utils.CheckLuhn(orderID); !ok {
		http.Error(resp, "invalid order ID", http.StatusUnprocessableEntity)
		return
	}

	if err := h.order.Add(ctx, userID, orderID); err != nil {
		if errors.Is(err, models.ErrOrderAlreadyUploadedByAnotherUser) {
			http.Error(resp, "order already uploaded by another user", http.StatusConflict)
			return
		}
		if errors.Is(err, models.ErrOrderInProcessed) {
			resp.WriteHeader(http.StatusOK)
			return
		}
		http.Error(resp, "internal error", http.StatusInternalServerError)
	}

	resp.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetAllOrders(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	ctx := req.Context()
	userID, err := h.getUserID(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusUnauthorized)
		return
	}

	orders, err := h.order.GetAll(ctx, userID)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	responseRaw, err := json.Marshal(orders)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(responseRaw)
}
