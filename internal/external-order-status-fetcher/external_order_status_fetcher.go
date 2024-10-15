package externalorderstatusfetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"go.uber.org/zap"
)

type ExternalOrderStatusFetcher struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func New(cfg *config.Config, logger *zap.SugaredLogger) *ExternalOrderStatusFetcher {
	return &ExternalOrderStatusFetcher{
		cfg:    cfg,
		logger: logger,
	}
}

func (e *ExternalOrderStatusFetcher) GetOrder(ctx context.Context, number string) (models.ExternalOrder, error) {
	var order models.ExternalOrder

	client := http.Client{}
	for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
		req, err := http.NewRequest(http.MethodGet, e.cfg.AccrualSystemAddress+"/api/orders/"+number, nil)
		if err != nil {
			return order, err
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			return order, err
		}
		if res.StatusCode == http.StatusTooManyRequests {
			time.Sleep(delay)
			continue
		}
		rawData, err := io.ReadAll(res.Body)
		defer res.Body.Close()

		if err != nil {
			e.logger.Errorln("error parse", err)
			return order, err
		}

		if len(rawData) == 0 {
			return order, errors.New("not order")
		}

		err = json.Unmarshal(rawData, &order)
		if err != nil {
			e.logger.Errorln("error unmarshal", err)
			return order, err
		}
		return order, nil

	}
	return order, fmt.Errorf("failed to fetch order after multiple retries")
}
