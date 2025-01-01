package manager_test

import (
	"errors"
	"os"
	"testing"

	"github.com/Isaac-Fate/myst/internal/manager"
	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/joho/godotenv"
)

func TestAddSecret(t *testing.T) {
	secretManager, err := createSecretManager()

	if err != nil {
		t.Fatal(err)
	}

	err = secretManager.AddSecret(&models.Secret{
		Key:            "test-secret",
		EncryptedValue: "xxx",
		Notes:          "this is a test secret",
	})

	if err != nil {
		t.Error(err)
	}

}

func TestFindSecrets(t *testing.T) {
	secretManager, err := createSecretManager()

	if err != nil {
		t.Fatal(err)
	}

	secrets, err := secretManager.FindSecrets("secret")

	if err != nil {
		t.Error(err)
	}

	for _, secret := range secrets {
		t.Log(secret)
	}
}

func createSecretManager() (*manager.SecretManager, error) {
	// Load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}

	secretStorePath := os.Getenv("TEST_SECRET_STORE_PATH")
	if secretStorePath == "" {
		return nil, errors.New("TEST_SECRET_STORE_PATH is not set")
	}

	indexPath := os.Getenv("TEST_SECRET_INDEX_PATH")
	if indexPath == "" {
		return nil, errors.New("TEST_SECRET_INDEX_PATH is not set")
	}

	secretManager, err := manager.NewSecretManager(secretStorePath, indexPath)
	if err != nil {
		return nil, err
	}

	return secretManager, nil
}
