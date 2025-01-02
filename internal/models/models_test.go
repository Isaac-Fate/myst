package models_test

import (
	"errors"
	"os"
	"testing"

	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInsert(t *testing.T) {
	// Get the test secret store database path
	secretStorePath, err := getTestSecretStorePath()
	if err != nil {
		t.Error(err)
	}

	// Open the database
	db, err := gorm.Open(sqlite.Open(secretStorePath), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	// Migrate the database
	db.AutoMigrate(&models.Secret{})

	// Create a new secret
	secret := models.Secret{
		Key:            "test-secret",
		EncryptedValue: "xxx",
		Notes:          "this is a test secret",
	}

	// Insert into the database
	err = db.Create(&secret).Error
	if err != nil {
		t.Error(err)
	}
}

func getTestSecretStorePath() (string, error) {
	// Load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		return "", err
	}

	// Get the test secret store database path
	secretStorePath := os.Getenv("TEST_SECRET_STORE_PATH")

	if secretStorePath == "" {
		return "", errors.New("TEST_SECRET_STORE_PATH is not set")
	}

	return secretStorePath, nil
}
