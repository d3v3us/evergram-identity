package main

import (
	"fmt"
	"log"
	"net"

	pbAuth "github.com/deveusss/evergram-identity/proto/auth"

	"github.com/deveusss/evergram-core/caching"
	"github.com/deveusss/evergram-core/config"
	"github.com/deveusss/evergram-identity/internal/account"
	"github.com/deveusss/evergram-identity/internal/auth"
	"github.com/deveusss/evergram-identity/internal/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config().Server.Port))
	if err != nil {
		panic(err)
	}
	database.ConnectDB()
	cache, _ := caching.NewAppCache()
	secret := config.Config().Auth.JwtSecret()
	accountRepository := account.NewAccountRepository(database.DB)
	authService := auth.NewAuthService(accountRepository, cache, secret)
	s := grpc.NewServer()
	reflection.Register(s)
	pbAuth.RegisterAuthServiceServer(s, authService)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
