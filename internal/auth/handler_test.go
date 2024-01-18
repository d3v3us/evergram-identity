package auth

import (
	"github.com/deveusss/evergram-core/caching"
	common "github.com/deveusss/evergram-core/common"
	"github.com/deveusss/evergram-core/encryption"
	"github.com/deveusss/evergram-core/validation"
	account "github.com/deveusss/evergram-identity/internal/account"
	"github.com/deveusss/evergram-identity/internal/account/registration"

	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	bcrypt "golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Connect to the SQLite in-memory database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory database")
	}

	// Migrate the database schema
	err = db.AutoMigrate(&account.Account{})
	if err != nil {
		panic("failed to migrate database schema")
	}

	return db
}

func generateBcryptHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
func TestLoginHandler_With_Valid_Credentials(t *testing.T) {
	// Set up Fiber app
	app := fiber.New()

	db := setupTestDB()
	cache, _ := caching.NewAppCache()
	secret := encryption.NewSecureString("secretKey")
	accountRepository := account.NewAccountRepository(db)
	registrationService := registration.NewRegistrationService(accountRepository, cache, secret)

	// Create the AuthService with the AccountRepository
	authService := NewAuthService(accountRepository, cache, secret)

	// Create the AuthHandler with the AuthService
	authHandler := NewAuthHandler(authService, registrationService)

	app.Post("/login", authHandler.Login)

	// Create a test user account and save it to the database
	testUser := &account.Account{
		Name:         "testuser1",
		Email:        "testuser1@example.com",
		PasswordHash: generateBcryptHash("password123"), // bcrypt hash for "password123"
	}

	err := accountRepository.Create(testUser)
	assert.NoError(t, err)

	// Test case: valid credentials
	reqBody := []byte(`{"nameoremail": "testuser1", "password": "password123"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response body
	var authResponse common.SucceededResponse[AuthResult]
	json.NewDecoder(resp.Body).Decode(&authResponse)

	// Make assertions on the response
	assert.NotEmpty(t, authResponse.Result.Token)
}

func TestLoginHandler_With_Invalid_Credentials(t *testing.T) {
	// Set up Fiber app
	app := fiber.New()

	db := setupTestDB()
	cache, _ := caching.NewAppCache()
	secret := encryption.NewSecureString("secretKey")

	// Create the AccountRepository with the in-memory database
	accountRepository := account.NewAccountRepository(db)
	registrationService := registration.NewRegistrationService(accountRepository, cache, secret)

	// Create the AuthService with the AccountRepository
	authService := NewAuthService(accountRepository, cache, secret)

	// Create the AuthHandler with the AuthService
	authHandler := NewAuthHandler(authService, registrationService)

	app.Post("/login", authHandler.Login)

	// Test case: invalid credentials
	reqBody := []byte(`{"nameoremail": "invaliduser", "password": "wrongpassword"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Parse the response body
	var errorResponse common.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	// Make assertions on the error response
	assert.Equal(t, errorResponse.Err, ErrUserWithSpecifiedCredentialsNotFound.Error())
	assert.NotNil(t, errorResponse.Message)
}
func TestLoginHandler_For_Validation(t *testing.T) {
	// Set up Fiber app
	app := fiber.New()

	db := setupTestDB()
	cache, _ := caching.NewAppCache()
	secret := encryption.NewSecureString("secretKey")
	accountRepository := account.NewAccountRepository(db)
	registrationService := registration.NewRegistrationService(accountRepository, cache, secret)

	// Create the AuthService with the AccountRepository
	authService := NewAuthService(accountRepository, cache, secret)

	// Create the AuthHandler with the AuthService
	authHandler := NewAuthHandler(authService, registrationService)
	app.Post("/login", validation.ValidationMiddleware, authHandler.Login)

	// Test case: invalid credentials
	reqBody := []byte(`{"nameoremail": "invaliduser", "password": ""}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Parse the response body
	var errorResponse common.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	// Make assertions on the error response
	assert.Equal(t, "Validation error", errorResponse.Message)
	assert.NotNil(t, errorResponse.Message)
}
