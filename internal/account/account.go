package account

import (
	"strings"

	"time"

	"github.com/deveusss/evergram-core/encryption"
	pbAuth "github.com/deveusss/evergram-identity/proto/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type AccountRole string

const (
	RoleAdmin     AccountRole = "admin"
	RoleAccount   AccountRole = "account"
	RoleGuest     AccountRole = "guest"
	RoleModerator AccountRole = "moderator"
)

// User struct
type Account struct {
	ID                  uuid.UUID   `gorm:"type:uuid;primary_key;"`
	Name                string      `gorm:"uniqueIndex;not null" json:"name"`
	FirstName           string      `gorm:"null" json:"first_name"`
	LastName            string      `gorm:"null" json:"last_name"`
	Email               string      `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash        string      `gorm:"not null"`
	IsLoggedIn          bool        `gorm:"not null" json:"is_logged_in"`
	IsVerified          bool        `gorm:"not null" json:"is_verified"`
	IsDeleted           bool        `gorm:"not null" json:"is_deleted"`
	IsLocked            bool        `gorm:"not null" json:"is_locked"`
	IsBanned            bool        `gorm:"not null" json:"is_banned"`
	ResetPasswordToken  string      `gorm:"null" json:"reset_password_token"`
	ResetPasswordExpiry time.Time   `gorm:"null" json:"reset_password_expiry"`
	Role                AccountRole `gorm:"not null" json:"role"`
	ProfilePic          []byte      `gorm:"null" json:"profile_pic"`
	Followers           int         `gorm:"null" json:"followers"`
	Following           int         `gorm:"null" json:"following"`
	Views               int         `gorm:"null" json:"views"`
	Rating              float64     `gorm:"null" json:"rating"`
	Photos              int         `gorm:"null" json:"photos"`
	TwitterURL          string      `gorm:"null" json:"twitter_url"`
	InstagramURL        string      `gorm:"null" json:"instagram_url"`
	LinkedInURL         string      `gorm:"null" json:"linked_in_url"`
	YouTubeURL          string      `gorm:"null" json:"you_tube_url"`
	WebsiteURL          string      `gorm:"null" json:"website_url"`
	CreatedAt           time.Time   `gorm:"not null"`
	UpdatedAt           time.Time   `gorm:"null"`
	DeletedAt           time.Time   `gorm:"null"`
}

// Function for creation of new account with required fields
func NewAccount(name, email, passwordHash string, role AccountRole) *Account {
	return &Account{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		IsLoggedIn:   false,
		IsVerified:   false,
		IsDeleted:    false,
		IsLocked:     false,
		IsBanned:     false,
		Role:         role,
		CreatedAt:    time.Now(),
	}
}
func NewAccountId() uuid.UUID {
	return uuid.New()
}

func ExtractUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		// Email is not in the expected format
		return ""
	}
	return parts[0]
}

func (user *Account) BeforeCreate(*gorm.DB) error {
	user.ID = NewAccountId()

	return nil
}
func GenerateToken(name string, email string, accountId uuid.UUID, secret encryption.ISecureString, exp time.Duration) (string, *pbAuth.TokenClaims, error) {
	tokenClaims := &pbAuth.TokenClaims{
		AccountId: accountId.String(),
		Name:      name,
		Email:     email,
		Exp:       timestamppb.New(time.Now().Add(time.Hour * exp)),
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"accountId": tokenClaims.AccountId,
		"username":  tokenClaims.Name,
		"email":     tokenClaims.Email,
		"exp":       tokenClaims.Exp.AsTime().Unix(),
	})

	// Generate the token string
	tokenString, err := token.SignedString(secret.Get())
	if err != nil {
		return "", nil, err
	}

	// Return the token string and error
	return tokenString, tokenClaims, nil
}
