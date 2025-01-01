package manager

import (
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
