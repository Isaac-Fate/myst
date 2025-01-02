package handlers

import (
	"fmt"

	"github.com/Isaac-Fate/myst/cmd/context"
	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
	"github.com/atotto/clipboard"
	"github.com/manifoldco/promptui"
)

func FindSecrets(appContext *context.AppContext) error {
	// Prompt for search term
	prompt := promptui.Prompt{
		Label: "Enter search term",
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("search term cannot be empty")
			}
			return nil
		},
	}

	term, err := prompt.Run()
	if err != nil {
		return err
	}

	// Search for secrets
	secrets, err := appContext.SecretManager.FindSecrets(term)
	if err != nil {
		return fmt.Errorf("failed to search secrets: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found")
		return nil
	}

	// Display found secrets
	fmt.Printf("Found %d secrets:\n", len(secrets))
	for _, secret := range secrets {
		// Create a template for displaying secret details
		fmt.Printf("\n🔑 Key: %s\n", secret.Key)
		if secret.Website != "" {
			fmt.Printf("🌐 Website: %s\n", secret.Website)
		}
		if secret.Notes != "" {
			fmt.Printf("📝 Notes: %s\n", secret.Notes)
		}

		// Create a selection prompt for value actions
		actionPrompt := promptui.Select{
			Label: "Choose action for secret value",
			Items: []string{
				"Skip",
				"Display in terminal",
				"Copy to clipboard",
			},
		}

		idx, _, err := actionPrompt.Run()
		if err != nil {
			return err
		}

		if idx == 0 { // Skip
			continue
		}

		// Decrypt the secret value
		decryptedValue, err := mycrypto.Decrypt(appContext.Passphrase, secret.EncryptedValue)
		if err != nil {
			return fmt.Errorf("failed to decrypt secret value: %w", err)
		}

		switch idx {
		case 1: // Display in terminal
			fmt.Printf("🔒 Value: %s\n", decryptedValue)
		case 2: // Copy to clipboard
			if err := clipboard.WriteAll(decryptedValue); err != nil {
				return fmt.Errorf("failed to copy to clipboard: %w", err)
			}
			fmt.Printf("✅ Value copied to clipboard\n")
		}
	}

	return nil
}
