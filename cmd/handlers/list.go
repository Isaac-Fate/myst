package handlers

import (
	"fmt"

	"github.com/Isaac-Fate/myst/cmd/context"
	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
	"github.com/atotto/clipboard"
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

	// Ask if user wants to view/copy any secret values
	viewPrompt := promptui.Prompt{
		Label:     "View or copy any secret values",
		IsConfirm: true,
	}

	if result, err := viewPrompt.Run(); err == nil && result == "y" {
		// Create items for selection
		items := make([]string, len(secrets))
		for i, s := range secrets {
			items[i] = fmt.Sprintf("[%d] %s (%s)", i+1, s.Key, s.Notes)
		}

		// Let user select which secret to view/copy
		selectPrompt := promptui.Select{
			Label: "Select a secret",
			Items: items,
			Size:  10, // Show 10 items at a time if list is long
		}

		idx, _, err := selectPrompt.Run()
		if err != nil {
			return err
		}

		selectedSecret := secrets[idx]

		// Ask what to do with the selected secret
		actionPrompt := promptui.Select{
			Label: "Choose action",
			Items: []string{
				"Skip",
				"Display in terminal",
				"Copy to clipboard",
			},
		}

		actionIdx, _, err := actionPrompt.Run()
		if err != nil {
			return err
		}

		if actionIdx == 0 { // Skip
			return nil
		}

		// Decrypt the secret value
		decryptedValue, err := mycrypto.Decrypt(appContext.Passphrase, selectedSecret.EncryptedValue)
		if err != nil {
			return fmt.Errorf("failed to decrypt secret value: %w", err)
		}

		switch actionIdx {
		case 1: // Display in terminal
			fmt.Printf("\n🔒 Value for '%s': %s\n", selectedSecret.Key, decryptedValue)
		case 2: // Copy to clipboard
			if err := clipboard.WriteAll(decryptedValue); err != nil {
				return fmt.Errorf("failed to copy to clipboard: %w", err)
			}
			fmt.Printf("\n✅ Value for '%s' copied to clipboard\n", selectedSecret.Key)
		}
	}

	return nil
}
