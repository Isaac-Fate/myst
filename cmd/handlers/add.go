package handlers

import (
	"fmt"

	"github.com/Isaac-Fate/myst/cmd/context"
	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
	"github.com/Isaac-Fate/myst/internal/models"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
)

func AddSecret(appContext *context.AppContext) error {

	// Prompt for secret key
	prompt := promptui.Prompt{
		Label: "Enter the secret key",
	}

	secretKey, err := prompt.Run()
	if err != nil {
		return err
	}

	// Prompt for secret value
	prompt = promptui.Prompt{
		Label: "Enter the secret value",
		Mask:  '*',
	}

	value, err := prompt.Run()
	if err != nil {
		return err
	}

	// Encrypt the secret value
	encryptedValue, err := mycrypto.Encrypt(appContext.Passphrase, value)
	if err != nil {
		return err
	}

	// Prompt for website
	prompt = promptui.Prompt{
		Label: "Enter the website",
	}

	website, err := prompt.Run()
	if err != nil {
		return err
	}

	// Prompt for notes
	prompt = promptui.Prompt{
		Label: "Enter any notes",
	}

	notes, err := prompt.Run()
	if err != nil {
		return err
	}

	// Create the secret
	secret := models.Secret{
		ID:             uuid.New(),
		Key:            secretKey,
		EncryptedValue: encryptedValue,
		Website:        website,
		Notes:          notes,
	}

	// Add the secret
	if err := appContext.SecretManager.AddSecret(&secret); err != nil {
		return err
	}

	fmt.Println("✅ Secret added successfully!")
	return nil
}
