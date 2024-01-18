package auth

import (
	pb "github.com/deveusss/evergram-identity/proto/auth"
)

func Succeeded(token string, claims *pb.TokenClaims) *pb.AuthResponse {
	return &pb.AuthResponse{
		Token:         token,
		Authenticated: true,
		Claims:        claims,
	}
}
func Failed() *pb.AuthResponse {
	return &pb.AuthResponse{
		Token:         "",
		Authenticated: false,
		Claims:        nil,
	}
}
