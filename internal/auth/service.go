// service.go
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/deveusss/evergram-identity/internal/account"
	pbAuth "github.com/deveusss/evergram-identity/proto/auth"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/deveusss/evergram-core/encryption"

	"net/mail"

	"github.com/deveusss/evergram-core/caching"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	accountRepo *account.AccountRepository
	cache       *caching.AppCache
	secretKey   encryption.ISecureString // Secret key for token signing
	pbAuth.UnimplementedAuthServiceServer
}

func NewAuthService(accountRepo *account.AccountRepository, cache *caching.AppCache, secretKey encryption.ISecureString) *AuthService {
	return &AuthService{
		accountRepo: accountRepo,
		cache:       cache,
		secretKey:   secretKey,
	}
}

func (s *AuthService) AuthenticateUser(ctx context.Context, req *pbAuth.AuthRequest) (*pbAuth.AuthResponse, error) {
	nameOrEmail := req.Username
	password := req.Password

	var user *account.Account
	var err error
	if isEmail(nameOrEmail) {
		user, err = s.accountRepo.GetByEmail(nameOrEmail)
	} else {
		user, err = s.accountRepo.GetByUsername(nameOrEmail)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Failed(), err
	}

	// Check if the user exists and validate the password
	if user == nil {
		return Failed(), ErrUserWithSpecifiedCredentialsNotFound
	}
	if !CheckPasswordHash(password, user.PasswordHash) {
		return Failed(), ErrInvalidCredentials
	}

	// Generate and return the JWT token
	token, claims, err := account.GenerateToken(user.Name, user.Email, user.ID, s.secretKey)
	if err != nil {
		return Failed(), err
	}

	return Succeeded(token, claims), nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pbAuth.ValidateTokenRequest) (*pbAuth.TokenClaims, error) {
	claims := jwt.MapClaims{}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, ErrInvalidJwtTokenSigningMethod
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, ErrInvalidJwtToken
	}

	return &pbAuth.TokenClaims{
		AccountId: claims["accountId"].(string),
		Name:      claims["username"].(string),
		Email:     claims["email"].(string),
		Exp:       timestamppb.New(time.Unix(int64(claims["exp"].(float64)), 0)),
	}, nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
