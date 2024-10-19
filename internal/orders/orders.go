package orders

import (
	"context"
	"fmt"
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
	processCtx                 context.Context
	store                      Storer
	cfg                        *config.Config
	wg                         *sync.WaitGroup
	externalOrderStatusFetcher ExternalOrderStatusFetcher
	logger                     *zap.SugaredLogger
}

func New(
	processCtx context.Context,
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
		processCtx: processCtx,
	}
}

func (s *Service) Add(ctx context.Context, userID uint64, orderID string) error {
	order, err := s.store.AddNewOrder(ctx, userID, orderID)
	if err != nil {
		return err
	}

	go func() {
		s.updateOrderInBackground(order)
	}()

	return nil
}

func (s *Service) GetAll(ctx context.Context, userID uint64) ([]models.Order, error) {
	return s.store.GetAllOrders(ctx, userID)
}

func (s *Service) updateOrderInBackground(order models.Order) {
	delay := 2 * time.Second
	maxDelay := 1 * time.Minute

	for {
		select {
		case <-s.processCtx.Done():
			s.Shutdown()
			return
		default:
			orderStatus, err := s.externalOrderStatusFetcher.GetOrder(s.processCtx, order.Number)
			if err != nil {
				s.logger.Errorln("Failed to fetch order status:", err)
				return
			}

			s.logger.Infoln("Fetched order status:", orderStatus.Status)

			if orderStatus.Status == string(models.PROCESSED) || orderStatus.Status == string(models.INVALID) {
				order.Accrual = orderStatus.Accrual
				order.Status = models.OrderStatus(orderStatus.Status)
				err = s.store.UpdateOrder(s.processCtx, order)
				if err != nil {
					s.logger.Errorln("Failed to update order:", err)
				} else {
					s.logger.Infoln("Successfully updated order:", order.Number, order.Accrual)
				}
				return
			}

			time.Sleep(delay)
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}
}

func (s *Service) Shutdown() {
	s.wg.Wait()
}
