package core

import (
	"ecommerce/app/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	//DSN and Connecting to the db
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPassword, dbName, dbPort)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return err
	}
	log.Println("Connected to database successfully")

	// Define models for auto migration
	Models := []interface{}{
		&models.User{},
		&models.BlacklistedToken{},
		&models.EmailVerificationToken{},
		&models.PasswordResetToken{},
		&models.Addon{},
		&models.Branch{},
		&models.Product{},
		&models.ShippingAddress{},
		&models.VariationOption{},
		&models.ProductVariation{},
		&models.Category{},
		&models.Notification{},
		&models.SubCategory{},
		&models.Review{},
		&models.OrderItemAddon{},
		&models.OrderItem{},
		&models.OrderItemVariation{},
		&models.Payment{},
	}

	// Loop for each model for auto Migration
	for _, model := range Models {
		err := DB.AutoMigrate(model)
		if err != nil {
			log.Fatalf("failed to migrate database: %v", err)
			return err
		}
		log.Printf("Migrated model %T successfully", model)
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
