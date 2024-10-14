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
	order, err := s.store.AddNewOrder(ctx, userID, orderID)
	if err != nil {
		return err
	}

	go func() {
		s.updateOrderInBackground(context.Background(), order)
	}()

	return nil
}

func (s *Service) GetAll(ctx context.Context, userID uint64) ([]models.Order, error) {
	return s.store.GetAllOrders(ctx, userID)
}

func (s *Service) updateOrderInBackground(ctx context.Context, order models.Order) {
	fmt.Println(s.logger)
	s.logger.Infoln("Start updating order:", order.Number, order.Accrual)

	for {
		orderStatus, err := s.externalOrderStatusFetcher.GetOrder(ctx, order.Number)
		if err != nil {
			s.logger.Errorln("Failed to fetch order status:", err)
			return
		}

		s.logger.Infoln("Fetched order status:", orderStatus.Status)

		if orderStatus.Status == string(models.PROCESSED) || orderStatus.Status == string(models.INVALID) {
			order.Accrual = orderStatus.Accrual
			order.Status = models.OrderStatus(orderStatus.Status)
			err = s.store.UpdateOrder(ctx, order)
			if err != nil {
				s.logger.Errorln("Failed to update order:", err)
			} else {
				s.logger.Infoln("Successfully updated order:", order.Number, order.Accrual)
			}
			break
		}

		time.Sleep(2 * time.Second)
	}
}

// func (s *Service) CheckStatus(ctx context.Context) {
// 	timer := time.NewTicker(2 * time.Second)
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-timer.C:
// 			orders, err := s.store.GetPendingOrders(ctx)
// 			if err != nil {
// 				s.logger.Errorln("get order error", err)
// 				return
// 			}

// 			jobs := make(chan models.Order)
// 			go func() {
// 				for _, order := range orders {
// 					jobs <- order
// 				}
// 				close(jobs)
// 			}()

// 			for i := 1; i <= 3; i++ {
// 				s.wg.Add(1)
// 				go s.processing(ctx, jobs)
// 			}

// 			s.wg.Wait()
// 		}
// 	}
// }

// func (s *Service) processing(ctx context.Context, jobs chan models.Order) {
// 	defer s.wg.Done()
// 	for job := range jobs {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			orderStatus, err := s.externalOrderStatusFetcher.GetOrder(ctx, job.Number)
// 			if err != nil {
// 				s.logger.Errorln("get order status error", err)
// 				continue
// 			}

// 			if orderStatus.Status != string(job.Status) {
// 				job.Accrual = orderStatus.Accrual
// 				job.Status = models.OrderStatus(orderStatus.Status)
// 				err = s.store.UpdateOrder(ctx, job)
// 				if err != nil {
// 					s.logger.Errorln("update order error", err)
// 					continue // Продолжим обработку других заказов
// 				}
// 			}
// 		}
// 	}
// }
