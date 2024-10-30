package orders

import (
	"context"
	"sync"
	"time"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"go.uber.org/zap"
)

type Storer interface {
	AddNewOrder(ctx context.Context, userID uint64, orderID string) (models.Order, error)
	GetAllOrders(ctx context.Context, userID uint64) ([]models.Order, error)
	GetPendingOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, orderID string) (models.Order, error)
}

type ExternalOrderStatusFetcher interface {
	GetOrder(ctx context.Context, orderNumber string) (models.ExternalOrder, error)
}

type Service struct {
	store                      Storer
	cfg                        *config.Config
	wg                         *sync.WaitGroup
	externalOrderStatusFetcher ExternalOrderStatusFetcher
	logger                     *zap.SugaredLogger
}

func New(
	store Storer,
	externalOrderStatusFetcher ExternalOrderStatusFetcher,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		store:                      store,
		cfg:                        cfg,
		externalOrderStatusFetcher: externalOrderStatusFetcher,
		wg:                         &sync.WaitGroup{},
		logger:                     logger,
	}
}

func (s *Service) Add(ctx context.Context, userID uint64, orderID string) error {
	_, err := s.store.AddNewOrder(ctx, userID, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAll(ctx context.Context, userID uint64) ([]models.Order, error) {
	return s.store.GetAllOrders(ctx, userID)
}

func (s *Service) UpdateOrder(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.wg.Wait()
			return
		case <-ticker.C:
			s.processPendingOrders(ctx)
		}
	}
}

func (s *Service) processPendingOrders(ctx context.Context) {
	orders, err := s.store.GetPendingOrders(ctx)
	if err != nil {
		s.logger.Errorln("Failed to get pending orders:", err)
		return
	}

	if len(orders) == 0 {
		s.logger.Infoln("No pending orders to process")
		return
	}

	jobs := make(chan models.Order, len(orders))
	go func() {
		for _, order := range orders {
			jobs <- order
		}
		close(jobs)
	}()

	for i := 0; i < s.cfg.NumWorkers; i++ {
		s.wg.Add(1)
		go s.worker(ctx, jobs)
	}
}

func (s *Service) worker(ctx context.Context, jobs chan models.Order) {
	defer s.wg.Done()

	for order := range jobs {
		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		orderStatus, err := s.externalOrderStatusFetcher.GetOrder(ctxTimeout, order.Number)
		cancel()

		if err != nil {
			s.logger.Errorf("Failed to get status for order %s: %v", order.Number, err)
			continue
		}

		order.Status = models.OrderStatus(orderStatus.Status)
		order.Accrual = orderStatus.Accrual

		if err := s.store.UpdateOrder(ctx, order); err != nil {
			s.logger.Errorf("Failed to update order %s: %v", order.Number, err)
		}
	}
}
