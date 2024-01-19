package main

import (
	"fmt"
	"log"
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

func main() {
	cfg := config.Load[config.AppConfig]()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Config.GRPC.Port))
	if err != nil {
		panic(err)
	}
	cache, _ := caching.NewAppCache()
	db, err := database.New(&cfg.Config.DbConfig)
	accountRepository := account.NewAccountRepository(db)
	authService := auth.NewAuthService(accountRepository, cache, cfg.Config)
	s := grpc.NewServer()
	reflection.Register(s)
	pbAuth.RegisterAuthServiceServer(s, authService)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
