package app

import (
	"context"
	"database/sql"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/app/server"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/balances"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/external-order-status-fetcher"
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
	ctx    context.Context
	server *server.Server
	db     *sql.DB
	logger *zap.SugaredLogger
	order  *orders.Service
}

func New(ctx context.Context) *App {
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
	externalOrderStatusFetcher := externalorderstatusfetcher.New(cfg, logger)
	orderService := orders.New(ctx, store, externalOrderStatusFetcher, cfg, logger)
	balanceService := balances.New(store)

	handler := handler.New(userService, orderService, balanceService, cfg)

	server := server.New(*cfg, logger, handler, orderService, tokenService)

	return &App{
		server: server,
		db:     db,
		logger: logger,
		order:  orderService,
		ctx: ctx,
	}
}

func (a *App) RunContext() {
	a.logger.Infow("Running migrations", "db", a.db)
	if err := goose.UpContext(a.ctx, a.db, "./internal/store/migrations"); err != nil {
		a.logger.Fatalw("Failed to run migrations", "error", err)
	}

	a.server.Run(a.ctx)

	defer a.db.Close()
}
