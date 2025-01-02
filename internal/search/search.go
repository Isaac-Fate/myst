package search

import (
	"errors"

	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/blevesearch/bleve/v2"
	"gorm.io/gorm"
)

func OpenIndex(indexPath string) (bleve.Index, error) {
	// Open the index
	index, err := bleve.Open(indexPath)

	if err == nil {
		return index, nil
	}

	// Create a new index
	if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
		// Create a new index mapping
		mapping := bleve.NewIndexMapping()

		// Create a new index
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, err
		}

		return index, nil
	}

	return nil, err
}

func FindSecrets(db *gorm.DB, index bleve.Index, query string) ([]models.Secret, error) {
	// Find secret IDs
	secretIds, err := FindSecretIds(index, query)
	if err != nil {
		return nil, err
	}

	// Find secrets in the database by IDs
	var secrets []models.Secret
	err = db.Where("id IN ?", secretIds).Find(&secrets).Error
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func FindSecretIds(index bleve.Index, query string) ([]string, error) {
	// Create a match query
	searchQuery := bleve.NewMatchQuery(query)

	// Create a search request based on the query
	searchRequest := bleve.NewSearchRequest(searchQuery)

	// Search the index
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// Secret IDs
	secretIds := make([]string, searchResult.Hits.Len())

	// Collect each secret ID
	for _, hit := range searchResult.Hits {
		secretIds = append(secretIds, hit.ID)
	}

	return secretIds, nil
}
