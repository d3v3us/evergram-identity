package database

import (
	"fmt"

	"github.com/deveusss/evergram-identity/internal/account"

	"github.com/deveusss/evergram-core/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB represents the GORM database connection
var DB *gorm.DB

// ConnectDB connects to the PostgreSQL database
func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Config().Database.Host,
		config.Config().Database.Port,
		config.Config().Database.User,
		config.Config().Database.Password,
		config.Config().Database.Name)

	fmt.Println("Connecting to Database ...")

	fmt.Println(dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to")
	DB.AutoMigrate(&account.Account{})
	fmt.Println("Database Migrated")
}
