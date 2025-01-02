package models

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Key            string    `gorm:"unique"`
	EncryptedValue string    `gorm:"not null"`
	Website        string
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (secret *Secret) OmitEncryptedValue() Secret {
	return Secret{
		ID:        secret.ID,
		Key:       secret.Key,
		Website:   secret.Website,
		Notes:     secret.Notes,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}
}
