package handlers

import (
	"fmt"

	"github.com/Isaac-Fate/myst/cmd/context"
	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
	"github.com/manifoldco/promptui"
)

func UpdateSecret(appContext *context.AppContext) error {
	// First list all secrets for selection
	secrets, err := appContext.SecretManager.ListSecrets()
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found")
		return nil
	}

	// Create items for selection
	items := make([]string, len(secrets))
	for i, s := range secrets {
		desc := s.Key
		if s.Notes != "" {
			desc += fmt.Sprintf(" (%s)", s.Notes)
		}
		items[i] = fmt.Sprintf("[%d] %s", i+1, desc)
	}

	// Let user select which secret to update
	selectPrompt := promptui.Select{
		Label: "Select a secret to update",
		Items: items,
		Size:  10, // Show 10 items at a time if list is long
	}

	idx, _, err := selectPrompt.Run()
	if err != nil {
		return err
	}

	selectedSecret := secrets[idx]

	// Ask what to update
	updatePrompt := promptui.Select{
		Label: "What would you like to update",
		Items: []string{
			"Value",
			"Website",
			"Notes",
		},
	}

	updateIdx, _, err := updatePrompt.Run()
	if err != nil {
		return err
	}

	switch updateIdx {
	case 0: // Update value
		prompt := promptui.Prompt{
			Label: "Enter new value",
			Mask:  '*',
			Validate: func(input string) error {
				if len(input) == 0 {
					return fmt.Errorf("value cannot be empty")
				}
				return nil
			},
		}

		newValue, err := prompt.Run()
		if err != nil {
			return err
		}

		// Encrypt the new value
		encryptedValue, err := mycrypto.Encrypt(appContext.Passphrase, newValue)
		if err != nil {
			return err
		}

		selectedSecret.EncryptedValue = encryptedValue

	case 1: // Update website
		prompt := promptui.Prompt{
			Label:   "Enter new website",
			Default: selectedSecret.Website,
		}

		newWebsite, err := prompt.Run()
		if err != nil {
			return err
		}

		selectedSecret.Website = newWebsite

	case 2: // Update notes
		prompt := promptui.Prompt{
			Label:   "Enter new notes",
			Default: selectedSecret.Notes,
		}

		newNotes, err := prompt.Run()
		if err != nil {
			return err
		}

		selectedSecret.Notes = newNotes
	}

	// Confirm update
	confirmPrompt := promptui.Prompt{
		Label:     "Confirm update",
		IsConfirm: true,
	}

	result, err := confirmPrompt.Run()
	if err != nil {
		return nil // User cancelled
	}

	if result == "y" {
		if err := appContext.SecretManager.UpdateSecret(&selectedSecret); err != nil {
			return fmt.Errorf("failed to update secret: %w", err)
		}
		fmt.Printf("✅ Secret '%s' updated successfully!\n", selectedSecret.Key)
	}

	return nil
}
