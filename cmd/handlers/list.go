package handlers

import (
	"fmt"

	"github.com/Isaac-Fate/myst/cmd/context"
	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
	"github.com/manifoldco/promptui"
)

func ListSecrets(appContext *context.AppContext) error {
	secrets, err := appContext.SecretManager.ListSecrets()
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found")
		return nil
	}

	// Display all secrets
	fmt.Printf("Found %d secrets:\n", len(secrets))
	for i, secret := range secrets {
		fmt.Printf("\n[%d] 🔑 %s\n", i+1, secret.Key)
		if secret.Website != "" {
			fmt.Printf("    🌐 Website: %s\n", secret.Website)
		}
		if secret.Notes != "" {
			fmt.Printf("    📝 Notes: %s\n", secret.Notes)
		}
	}

	// Ask if user wants to view any secret values
	viewPrompt := promptui.Prompt{
		Label:     "View any secret values",
		IsConfirm: true,
	}

	if result, err := viewPrompt.Run(); err == nil && result == "y" {
		// Create a selection prompt for secrets
		items := make([]string, len(secrets))
		for i, s := range secrets {
			items[i] = fmt.Sprintf("%s (%s)", s.Key, s.Notes)
		}

		selectPrompt := promptui.Select{
			Label: "Select a secret to view",
			Items: items,
			Size:  10, // Show 10 items at a time if list is long
		}

		idx, _, err := selectPrompt.Run()
		if err != nil {
			return err
		}

		// Decrypt and show the selected secret value
		selectedSecret := secrets[idx]
		decryptedValue, err := mycrypto.Decrypt(appContext.Passphrase, selectedSecret.EncryptedValue)
		if err != nil {
			return fmt.Errorf("failed to decrypt secret value: %w", err)
		}
		fmt.Printf("\n🔒 Value for '%s': %s\n", selectedSecret.Key, decryptedValue)
	}

	return nil
}
