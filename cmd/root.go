/*
Copyright © 2024 Isaac Fei
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Isaac-Fate/myst/internal/config"
	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
	"github.com/Isaac-Fate/myst/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Configuration this program
var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "myst",
	Short: "MyST (My SecreTs) -- A Simple Secret Value Manager",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Print the help message
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myst.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Load the configuration
	// Prompt the user to set the configuration if it doesn't exist
	err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}

}

func getDataDir() (string, error) {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// The data directory of this program is ~/.myst
	dataDir := filepath.Join(homeDir, "myst")

	// Create the data directory if it doesn't exist
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return "", err
	}

	return dataDir, nil
}

func getConfigPath() (string, error) {
	// Get the data directory
	dataDir, err := getDataDir()
	if err != nil {
		return "", err
	}

	// The config file of this program is ~/myst/config.yml
	configPath := filepath.Join(dataDir, "config.yml")

	return configPath, nil
}

func loadConfig() error {
	// Get the config path
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Load the configuration, set the global variable cfg
	cfg, err = config.LoadConfig(configPath)

	if err == nil {
		return nil
	}

	// Guide the user to set the config
	if errors.Is(err, os.ErrNotExist) {
		return setConfig(configPath)
	}

	return nil
}

func setConfig(configPath string) error {
	// Prompt to set the filepath of the password store database
	setPasswordStorePathPrompt := promptui.Prompt{
		Label: "🔐 Enter the filepath of the password store",
		Validate: func(passwordStorePath string) error {
			// Resolve the path
			passwordStorePath, err := utils.ResolvePath(passwordStorePath)
			if err != nil {
				return err
			}

			// Check if the file exists
			_, err = os.Stat(passwordStorePath)
			if err == nil {
				return errors.New("the file at the specified path already exists")
			}

			// Check if the path ends with .db
			if !strings.HasSuffix(passwordStorePath, ".db") {
				return errors.New("password store path must end with .db")
			}

			return nil
		},
	}

	// Get the input path of the password store
	passwordStorePath, err := setPasswordStorePathPrompt.Run()

	if err != nil {
		return err
	}

	// Resolve the path
	passwordStorePath, err = utils.ResolvePath(passwordStorePath)
	if err != nil {
		return err
	}

	// Create the parent directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(passwordStorePath), 0755)
	if err != nil {
		return err
	}

	// Prompt to set the passphrase
	setPassphrasePrompt := promptui.Prompt{
		Label: "🔑 Enter the passphrase",
		Validate: func(passphrase string) error {
			// Check if the passphrase is at least 8 characters long
			if len(passphrase) < 8 {
				return errors.New("passphrase must be at least 8 characters long")
			}

			return nil
		},
		Mask: '*',
	}

	passphrase, err := setPassphrasePrompt.Run()

	if err != nil {
		return err
	}

	// Set the cfg
	cfg.SecretStorePath = passwordStorePath
	cfg.DigestedPassphrase = mycrypto.DigestPassphrase(passphrase)

	// Save the config
	err = cfg.Save(configPath)
	if err != nil {
		return err
	}

	return nil
}
