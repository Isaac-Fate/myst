package manager

import (
	"fmt"

	"github.com/Isaac-Fate/myst/internal/database"
	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/Isaac-Fate/myst/internal/search"
	"github.com/blevesearch/bleve/v2"
	"gorm.io/gorm"
)

type SecretManager struct {
	db    *gorm.DB
	index bleve.Index
}

func NewSecretManager(dbPath, indexPath string) (*SecretManager, error) {
	// Open the database
	db, err := database.OpenSecretStore(dbPath)

	if err != nil {
		return nil, err
	}

	// Open the index
	index, err := search.OpenIndex(indexPath)

	if err != nil {
		return nil, err
	}

	return &SecretManager{
		db:    db,
		index: index,
	}, nil
}

func (manager *SecretManager) AddSecret(secret *models.Secret) error {
	// Create a transaction
	tx := manager.db.Begin()

	// Add the secret to the database
	err := database.AddSecret(tx, secret)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Add the secret to the index
	err = search.AddSecret(manager.index, secret)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit
	tx.Commit()

	return nil
}

func (manager *SecretManager) FindSecrets(query string) ([]models.Secret, error) {
	// Create a transaction
	tx := manager.db.Begin()

	// Search the index and find the secret IDs
	secretIds, err := search.FindSecretIds(manager.index, query)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Find the secrets in the database
	secrets, err := database.GetSecrets(tx, secretIds)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit
	tx.Commit()

	return secrets, nil
}

// UpdateSecret updates an existing secret in both the database and search index
func (manager *SecretManager) UpdateSecret(secret *models.Secret) error {
	// Create a transaction
	tx := manager.db.Begin()

	// Update the secret in the database
	err := tx.Save(secret).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the secret in the search index
	err = search.UpdateSecret(manager.index, secret)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	tx.Commit()
	return nil
}

// RemoveSecret removes a secret from both the database and search index
func (manager *SecretManager) RemoveSecret(secret *models.Secret) error {
	// Create a transaction
	tx := manager.db.Begin()

	// Remove the secret from the database
	err := tx.Delete(secret).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove the secret from the search index
	err = search.RemoveSecret(manager.index, secret)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	tx.Commit()
	return nil
}

// GetSecret retrieves a secret by its ID
func (manager *SecretManager) GetSecret(id string) (*models.Secret, error) {
	return database.GetSecret(manager.db, id)
}

// ListSecrets returns all secrets in the database
func (manager *SecretManager) ListSecrets() ([]models.Secret, error) {
	var secrets []models.Secret
	err := manager.db.Find(&secrets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	return secrets, nil
}
