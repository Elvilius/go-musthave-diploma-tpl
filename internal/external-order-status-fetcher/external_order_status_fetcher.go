package externalorderstatusfetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"go.uber.org/zap"
)

type ExternalOrderStatusFetcher struct {
	cfg           *config.Config
	logger        *zap.SugaredLogger
	pauseDuration chan time.Duration
}

func New(cfg *config.Config, logger *zap.SugaredLogger) *ExternalOrderStatusFetcher {
	pauseDuration := make(chan time.Duration, 1)

	e := &ExternalOrderStatusFetcher{
		cfg:           cfg,
		logger:        logger,
		pauseDuration: pauseDuration,
	}
	go e.watchForPause()
	return e
}

func (e *ExternalOrderStatusFetcher) GetOrder(ctx context.Context, number string) (models.ExternalOrder, error) {
	var order models.ExternalOrder

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.cfg.AccrualSystemAddress+"/api/orders/"+number, nil)
	if err != nil {
		return order, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return order, err
	}
	if res.StatusCode == http.StatusTooManyRequests {
		retryAfter := res.Header.Get("Retry-After")
		sleepDuration, parseErr := e.parseRetryAfter(retryAfter)
		if parseErr != nil {
			e.logger.Errorln("failed to parse Retry-After header:", parseErr)
			return order, parseErr
		}
		e.pauseDuration <- sleepDuration
		return order, fmt.Errorf("retrying after %v", sleepDuration)
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

func (e *ExternalOrderStatusFetcher) parseRetryAfter(retryAfter string) (time.Duration, error) {
	seconds, err := strconv.Atoi(retryAfter)
	if err != nil {
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}

func (e *ExternalOrderStatusFetcher) watchForPause() {
	for {
		seconds := <-e.pauseDuration
		e.logger.Infof("Pausing all operations for %v", seconds)
		time.Sleep(seconds)
	}
}