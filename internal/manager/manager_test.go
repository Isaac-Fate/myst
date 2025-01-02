package manager_test

import (
	"errors"
	"os"
	"testing"

	"github.com/Isaac-Fate/myst/internal/manager"
	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestAddSecret(t *testing.T) {
	secretManager, err := createSecretManager()

	if err != nil {
		t.Fatal(err)
	}

	err = secretManager.AddSecret(&models.Secret{
		ID:             uuid.New(),
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

func TestUpdateSecret(t *testing.T) {
	secretManager, err := createSecretManager()
	if err != nil {
		t.Fatal(err)
	}

	// First create a secret
	secret := &models.Secret{
		ID:             uuid.New(),
		Key:            "update-test-secret",
		EncryptedValue: "original-value",
		Notes:          "original notes",
	}

	err = secretManager.AddSecret(secret)
	if err != nil {
		t.Fatal(err)
	}

	// Update the secret
	secret.Notes = "updated notes"
	secret.EncryptedValue = "new-value"

	err = secretManager.UpdateSecret(secret)
	if err != nil {
		t.Error(err)
	}

	// Verify the update
	updated, err := secretManager.GetSecret(secret.ID.String())
	if err != nil {
		t.Error(err)
	}

	if updated.Notes != "updated notes" {
		t.Errorf("expected updated notes, got %s", updated.Notes)
	}
}

func TestRemoveSecret(t *testing.T) {
	secretManager, err := createSecretManager()
	if err != nil {
		t.Fatal(err)
	}

	// First create a secret
	secret := &models.Secret{
		Key:            "remove-test-secret",
		EncryptedValue: "test-value",
		Notes:          "test notes",
	}

	err = secretManager.AddSecret(secret)
	if err != nil {
		t.Fatal(err)
	}

	// Remove the secret
	err = secretManager.RemoveSecret(secret)
	if err != nil {
		t.Error(err)
	}

	// Verify the removal
	_, err = secretManager.GetSecret(secret.ID.String())
	if err == nil {
		t.Error("expected error when getting removed secret")
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
