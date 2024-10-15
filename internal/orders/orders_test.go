package orders

import (
	"context"
	"errors"
	"testing"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	external_order_status_fetcher "github.com/Elvilius/go-musthave-diploma-tpl.git/internal/external-order-status-fetcher"
	mocks_orders "github.com/Elvilius/go-musthave-diploma-tpl.git/internal/mocks"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Add(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	type args struct {
		ctx     context.Context
		userID  uint64
		orderID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create new Order",
			fields: fields{
				cfg: &config.Config{},
			},
			args: args{
				ctx:     context.TODO(),
				userID:  1,
				orderID: "132123123123",
			},
		},
	}

	logger, _ := logger.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks_orders.NewMockOrderStore(ctrl)

			s := &Service{
				store:                      m,
				cfg:                        tt.fields.cfg,
				logger:                     logger,
				externalOrderStatusFetcher: &external_order_status_fetcher.ExternalOrderStatusFetcher{},
			}

			m.EXPECT().AddNewOrder(tt.args.ctx, tt.args.userID, tt.args.orderID).Return(models.Order{Number: "1234"}, nil)

			err := s.Add(tt.args.ctx, tt.args.userID, tt.args.orderID)

			require.NoError(t, err)
		})
	}
}

func TestService_GetAll(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	type args struct {
		ctx    context.Context
		userID uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockOrder []models.Order
		err       error
	}{
		{
			name: "Get all orders success",
			fields: fields{
				cfg: &config.Config{},
			},
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			mockOrder: []models.Order{
				{Number: "1", Accrual: 12},
				{Number: "1", Accrual: 12},
			},
			err: nil,
		},
		{
			name: "Get all orders err",
			fields: fields{
				cfg: &config.Config{},
			},
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			mockOrder: []models.Order{},
			err:       errors.New("err database"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks_orders.NewMockOrderStore(ctrl)

			s := &Service{
				store: m,
				cfg:   tt.fields.cfg,
			}

			m.EXPECT().GetAllOrders(tt.args.ctx, tt.args.userID).Return(tt.mockOrder, tt.err)

			orders, err := s.GetAll(tt.args.ctx, tt.args.userID)
			if err == nil {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.mockOrder, orders)
		})
	}
}
