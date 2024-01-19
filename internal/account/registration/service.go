package registration

import (
	"errors"

	"github.com/deveusss/evergram-identity/internal/account"

	"gorm.io/gorm"

	"github.com/deveusss/evergram-core/caching"
	"github.com/deveusss/evergram-core/config"
	"github.com/deveusss/evergram-core/encryption"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationService interface {
	Register(req *RegistrationRequest) (*RegistrationResult, error)
}
type RegistrationServiceImpl struct {
	accountRepo *account.AccountRepository
	cache       *caching.AppCache
	secretKey   encryption.ISecureString // Secret key for token signing
	config      *config.AppConfig
}

func NewRegistrationService(accountRepo *account.AccountRepository, cache *caching.AppCache,
	secretKey encryption.ISecureString,
	config *config.AppConfig) RegistrationService {
	return &RegistrationServiceImpl{
		accountRepo: accountRepo,
		cache:       cache,
		secretKey:   secretKey,
		config:      config,
	}
}
func (s *RegistrationServiceImpl) Register(req *RegistrationRequest) (*RegistrationResult, error) {
	var user *account.Account
	user, err := s.accountRepo.GetByEmail(req.Name)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Failed(), err
	}
	if user != nil {
		return Failed(), AccountAlreadyRegisteredError
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return Failed(), err
	}
	user = account.NewAccount(req.Name, req.Email, string(hash), account.RoleAccount)
	//Try to register user
	if err := s.accountRepo.Create(user); err != nil {
		return Failed(), err
	}
	token, claims, err := account.GenerateToken(user.Name, user.Email, user.ID, s.secretKey, s.config.AuthConfig.Jwt.TokenTTL)
	if err != nil {
		return Failed(), err
	}
	return Succeeded(token, claims), nil
}
