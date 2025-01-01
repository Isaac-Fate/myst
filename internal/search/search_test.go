package search_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/Isaac-Fate/myst/internal/search"
	"github.com/blevesearch/bleve/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestSearch(t *testing.T) {

	indexPath, err := getTestSeretIndexPath()
	if err != nil {
		t.Error(err)
	}

	index, err := openIndex(indexPath)
	if err != nil {
		t.Error(err)
	}

	secrets := []models.Secret{
		{
			ID:    uuid.New(),
			Key:   "test-secret-1",
			Notes: "github token",
		},
		{
			ID:    uuid.New(),
			Key:   "test-secret-2",
			Notes: "for deepseek service",
		},
		{
			ID:    uuid.New(),
			Key:   "test-secret-3",
			Notes: "google cloud",
		},
	}

	for _, secret := range secrets {
		err = index.Index(secret.ID.String(), &secret)
		if err != nil {
			t.Error(err)
		}
	}

	secretIds, err := search.FindSecretIds(index, "github")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(secretIds)

}

func openIndex(indexPath string) (bleve.Index, error) {
	index, err := bleve.Open(indexPath)
	if err == nil {
		return index, nil
	}

	// Try create a new index

	// Create a new index mapping
	mapping := bleve.NewIndexMapping()

	return bleve.New(indexPath, mapping)
}

func getTestSeretIndexPath() (string, error) {
	// Load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		return "", err
	}

	indexPath := os.Getenv("TEST_SECRET_INDEX_PATH")

	if indexPath == "" {
		return "", errors.New("TEST_SECRET_INDEX_PATH is not set")
	}

	return indexPath, nil
}
