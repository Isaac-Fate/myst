package handlers

import (
	"fmt"

	"github.com/Isaac-Fate/myst/cmd/context"
	"github.com/manifoldco/promptui"
)

func RemoveSecret(appContext *context.AppContext) error {
	// Prompt for the secret key
	prompt := promptui.Prompt{
		Label: "Enter the secret key to remove",
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("key cannot be empty")
			}
			return nil
		},
	}

	key, err := prompt.Run()
	if err != nil {
		return err
	}

	// Search for the secret
	secrets, err := appContext.SecretManager.FindSecrets(key)
	if err != nil {
		return fmt.Errorf("failed to find secret: %w", err)
	}

	if len(secrets) == 0 {
		return fmt.Errorf("secret with key '%s' not found", key)
	}

	// If multiple secrets found, let user select which one to remove
	var secretToRemove = secrets[0]
	if len(secrets) > 1 {
		items := make([]string, len(secrets))
		for i, s := range secrets {
			items[i] = fmt.Sprintf("%s (%s)", s.Key, s.Notes)
		}

		selectPrompt := promptui.Select{
			Label: "Multiple secrets found. Select one to remove",
			Items: items,
		}

		idx, _, err := selectPrompt.Run()
		if err != nil {
			return err
		}
		secretToRemove = secrets[idx]
	}

	// Confirm removal
	confirmPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure you want to remove secret '%s'", secretToRemove.Key),
		IsConfirm: true,
	}

	result, err := confirmPrompt.Run()
	if err != nil {
		return nil // User cancelled
	}

	if result == "y" {
		if err := appContext.SecretManager.RemoveSecret(&secretToRemove); err != nil {
			return fmt.Errorf("failed to remove secret: %w", err)
		}
		fmt.Printf("✅ Secret '%s' removed successfully\n", secretToRemove.Key)
	}

	return nil
}
