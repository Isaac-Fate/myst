/*
Copyright © 2024 Isaac Fei
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Isaac-Fate/myst/cmd/context"
	"github.com/Isaac-Fate/myst/cmd/handlers"
	"github.com/Isaac-Fate/myst/internal/config"
	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
	"github.com/Isaac-Fate/myst/internal/manager"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Global context for the application
var appContext = context.AppContext{}

var rootCmd = &cobra.Command{
	Use:   "myst",
	Short: "MyST (My SecreTs) -- A Simple Secret Value Manager",
	Long:  `A simple and secure secret manager for storing and retrieving sensitive information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// First, ensure we have a valid configuration
		if err := initializeConfig(); err != nil {
			return err
		}

		// Prompt for passphrase
		if err := loadPassphrase(); err != nil {
			return err
		}

		// Then initialize the secret manager
		if err := initializeSecretManager(); err != nil {
			return err
		}

		// Start the interactive command loop
		return startCommandLoop()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Initialize the configuration
func initializeConfig() error {
	// Try to load existing config
	err := config.LoadConfig(&appContext.Config)
	if err == nil {
		return nil
	}

	// If config doesn't exist, create it
	if errors.Is(err, os.ErrNotExist) {
		if err := createInitialConfig(); err != nil {
			return fmt.Errorf("failed to create initial configuration: %w", err)
		}
		return nil
	}

	return err
}

func createInitialConfig() error {
	fmt.Println("🔐 Welcome to MyST! Let's set up your secret store.")

	// Prompt for passphrase
	passphrasePrompt := promptui.Prompt{
		Label: "Enter your master passphrase (min 8 characters)",
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) < 8 {
				return errors.New("passphrase must be at least 8 characters")
			}
			return nil
		},
	}

	passphrase, err := passphrasePrompt.Run()
	if err != nil {
		return err
	}

	// Create and save the configuration
	appContext.Config = config.Config{
		DigestedPassphrase: mycrypto.DigestPassphrase(passphrase),
	}

	if err := config.Save(&appContext.Config); err != nil {
		return err
	}

	fmt.Println("✅ Configuration created successfully!")
	return nil
}

// Load the passphrase from the user
func loadPassphrase() error {
	passphrasePrompt := promptui.Prompt{
		Label: "Enter your master passphrase (min 8 characters)",
		Mask:  '*',
	}

	inputPassphrase, err := passphrasePrompt.Run()
	if err != nil {
		return err
	}

	// Verify the passphrase
	if !mycrypto.VerifyPassphrase(inputPassphrase, appContext.Config.DigestedPassphrase) {
		return errors.New("wrong passphrase")
	}

	// Set the passphrase
	appContext.Passphrase = inputPassphrase

	return nil
}

func initializeSecretManager() error {
	var err error
	appContext.SecretManager, err = manager.NewSecretManager(config.SecretStorePath(), config.SecretIndexPath())
	if err != nil {
		return fmt.Errorf("failed to initialize secret manager: %w", err)
	}
	return nil
}

func startCommandLoop() error {
	commands := []struct {
		Name        string
		Description string
		Handler     func(appContext *context.AppContext) error
	}{
		{"add", "Add a new secret", handlers.AddSecret},
		{"find", "Search for secrets", handlers.FindSecrets},
		{"list", "List all secrets", handlers.ListSecrets},
		{"remove", "Remove a secret", handlers.RemoveSecret},
		{"help", "Show help message", handlers.ShowHelp},
		{"quit", "Exit the application", nil},
	}

	// Create a map for quick command lookup
	commandToHandlerMap := make(map[string]func(appContext *context.AppContext) error)
	for _, cmd := range commands {
		commandToHandlerMap[cmd.Name] = cmd.Handler
	}

	for {
		prompt := promptui.Select{
			Label: "Select an action (or type a command)",
			Items: commands,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "▸ {{ .Name | cyan }} - {{ .Description }}",
				Inactive: "  {{ .Name | white }} - {{ .Description }}",
				Selected: "✔ {{ .Name | green }} - {{ .Description }}",
			},
			// Enable searching through commands
			Searcher: func(input string, index int) bool {
				cmd := commands[index]
				input = strings.ToLower(input)
				name := strings.ToLower(cmd.Name)
				return strings.HasPrefix(name, input)
			},
		}

		i, inputCommand, err := prompt.Run()

		if err != nil {
			if err == promptui.ErrInterrupt {
				return nil
			}
			return err
		}

		// Handle the typed command if it's a direct match
		if handler, exists := commandToHandlerMap[strings.ToLower(inputCommand)]; exists {
			if handler == nil { // quit command
				fmt.Println("👋 Goodbye!")
				return nil
			}
			if err := handler(&appContext); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			continue
		}

		// Handle the selected command
		if commands[i].Name == "quit" {
			fmt.Println("👋 Goodbye!")
			return nil
		}

		if err := commands[i].Handler(&appContext); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
