package search

import (
	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/blevesearch/bleve/v2"
)

// Adds a secret to the index.
//
// The secret is indexed with the fields name, website, and notes. The secret
// ID is used as the document ID.
//
// The function returns an error if the indexing fails.
func AddSecret(index bleve.Index, secret *models.Secret) error {
	return index.Index(secret.ID.String(), models.Secret{
		Key:     secret.Key,
		Website: secret.Website,
		Notes:   secret.Notes,
	})
}

// Updates a secret in the index.
//
// The function first deletes the old secret document by its ID, and then adds
// the new secret to the index. The function returns an error if the update
// fails.
func UpdateSecret(index bleve.Index, secret *models.Secret) error {
	// Delete the old secret document by its ID
	err := index.Delete(secret.ID.String())
	if err != nil {
		return err
	}

	// Add the new secret
	return AddSecret(index, secret)
}

// Removes a secret from the index.
//
// The secret ID is used to identify the document to delete.
//
// The function returns an error if the deletion fails.
func RemoveSecret(index bleve.Index, secret *models.Secret) error {
	return index.Delete(secret.ID.String())
}
