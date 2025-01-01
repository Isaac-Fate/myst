package database

import (
	"errors"

	"github.com/Isaac-Fate/myst/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenSecretStore(path string) (*gorm.DB, error) {
	// Open the database
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Migrate
	err = db.AutoMigrate(&models.Secret{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddSecret(db *gorm.DB, secret *models.Secret) error {
	err := db.Create(secret).Error

	if err != nil {
		return err
	}

	return nil
}

// Gets a secret from the database by its ID.
//
// If the secret does not exist or the retrieval fails, an error is returned.
func GetSecret(db *gorm.DB, id string) (*models.Secret, error) {
	// Create an empty secret
	var secret models.Secret

	// Get the secret from the database
	err := db.Where("id = ?", id).First(&secret).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("secret not found")
		}

		return nil, err
	}

	return &secret, nil
}

func GetSecrets(db *gorm.DB, ids []string) ([]models.Secret, error) {
	// Secrets to return
	var secrets []models.Secret

	// Get the secrets from the database
	err := db.Where("id IN ?", ids).Find(&secrets).Error

	if err != nil {
		return nil, err
	}

	return secrets, nil
}
