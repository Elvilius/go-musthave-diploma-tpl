package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/handler"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/jwt"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	cfg     config.Config
	handler *handler.Handler
	router  *chi.Mux
	logger  *zap.SugaredLogger
	token   *jwt.Jwt
}

func New(cfg config.Config, logger *zap.SugaredLogger, handler *handler.Handler, token *jwt.Jwt) *Server {
	s := &Server{
		cfg:     cfg,
		handler: handler,
		router:  chi.NewRouter(),
		token:   token,
		logger:  logger,
	}

	s.initRout()
	return s
}

func (s *Server) initRout() {
	s.router.Use(middleware.Logging(*s.logger))

	s.router.Post("/api/user/register", s.handler.RegisterUser)
	s.router.Post("/api/user/login", s.handler.LoginUser)

	s.router.Group(func(r chi.Router) {
		r.Use(middleware.ValidateJWT(s.token))
		r.Post("/api/user/orders", s.handler.AddNewOrder)
		r.Get("/api/user/orders", s.handler.GetAllOrders)
		r.Get("/api/user/balance", s.handler.GetBalance)
		r.Post("/api/user/balance/withdraw", s.handler.Withdraw)
		r.Get("/api/user/withdrawals", s.handler.GetWithdraw)
	})
}

func (s *Server) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:    s.cfg.RunAddress,
		Handler: s.router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Fatalf("Start service Error: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Infoln("Shutting down server...")
	case <-stop:
		s.logger.Infoln("Received stop signal, shutting down gracefully...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	s.logger.Infoln("Server exiting gracefully")
}
