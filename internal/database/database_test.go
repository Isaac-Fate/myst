package database_test

import (
	"errors"
	"os"
	"testing"

	"github.com/Isaac-Fate/myst/internal/database"
	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestCrud(t *testing.T) {
	// Get the test secret store path
	secretStorePath, err := getTestSecretStorePath()
	if err != nil {
		t.Error(err)
	}

	// Open the database
	db, err := database.OpenSecretStore(secretStorePath)
	if err != nil {
		t.Error(err)
	}

	// Create a secret
	secret := models.Secret{
		ID:             uuid.New(),
		Key:            "test-secret",
		EncryptedValue: "xxx",
		Notes:          "this is a test secret",
	}

	err = database.AddSecret(db, &secret)

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
