package app

import (
	"fmt"
	"log/slog"
	"net"

	pbAuth "github.com/deveusss/evergram-identity/proto/auth"

	"github.com/deveusss/evergram-core/caching"
	"github.com/deveusss/evergram-core/config"
	"github.com/deveusss/evergram-core/database"
	"github.com/deveusss/evergram-identity/internal/account"
	"github.com/deveusss/evergram-identity/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	cfg        *config.AppConfig
}

func NewApp(log *slog.Logger, cfg *config.AppConfig) *App {
	return &App{
		log:        log,
		port:       cfg.GRPC.Port,
		gRPCServer: nil,
		cfg:        cfg,
	}
}
func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}
func (a *App) run() error {
	a.log.Info("Starting identity server on port", slog.Int("port", a.cfg.GRPC.Port))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
	if err != nil {
		a.log.Error("Error listening: %v", err)
		panic(err)
	}

	cache, _ := caching.NewAppCache()
	db, err := database.New(a.log, &a.cfg.DbConfig)
	accountRepository := account.NewAccountRepository(db)
	authService := auth.NewAuthService(accountRepository, cache, a.cfg)
	gRPCServer := grpc.NewServer()
	reflection.Register(gRPCServer)
	pbAuth.RegisterAuthServiceServer(gRPCServer, authService)
	if err := gRPCServer.Serve(listener); err != nil {
		a.log.Error("Failed to serve: %v", err)
	}
	a.log.Info("Identity started")

	return nil

}
func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
	a.log.Info("Identity server stopped")

}
