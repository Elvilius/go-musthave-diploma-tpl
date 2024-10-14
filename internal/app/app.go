package app

import (
	"context"
	"database/sql"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/app/server"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/balances"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	external_order_status_fetcher "github.com/Elvilius/go-musthave-diploma-tpl.git/internal/external-order-status-fetcher"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/handler"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/orders"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/store"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/users"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/jwt"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type App struct {
	server *server.Server
	db     *sql.DB
	logger *zap.SugaredLogger
	order  *orders.Service
}

func New() *App {
	logger, err := logger.New()
	if err != nil {
		logger.Fatalw("Failed to open DB", "error", err)
	}
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DatabaseURI)
	if err != nil {
		logger.Fatalw("Failed to open DB", "error", err)
	}

	store := store.New(db)

	tokenService := jwt.New(cfg)
	userService := users.New(store, tokenService, cfg)
	externalOrderStatusFetcher := external_order_status_fetcher.New(cfg, logger)
	orderService := orders.New(store, externalOrderStatusFetcher, cfg, logger)
	balanceService := balances.New(store)

	handler := handler.New(userService, orderService, balanceService, cfg)

	server := server.New(*cfg, logger, handler, tokenService)

	return &App{
		server: server,
		db:     db,
		logger: logger,
		order:  orderService,
	}
}

func (a *App) RunContext(ctx context.Context) {
	a.logger.Infow("Running migrations", "db", a.db)
	if err := goose.UpContext(ctx, a.db, "./internal/store/migrations"); err != nil {
		a.logger.Fatalw("Failed to run migrations", "error", err)
	}

	// go a.order.CheckStatus(ctx)
	a.server.Run(ctx)

	defer a.db.Close()
}
