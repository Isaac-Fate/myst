package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const dataDirName = "myst"

type Config struct {
	DigestedPassphrase string `yaml:"digested_passphrase"`
}

func DataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dataDir := filepath.Join(homeDir, dataDirName)

	// Create the data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic(err)
	}

	return dataDir
}

func ConfigPath() string {
	return filepath.Join(DataDir(), "config.yml")
}

func SecretStorePath() string {
	return filepath.Join(DataDir(), "secret-store.db")
}

func SecretIndexPath() string {
	return filepath.Join(DataDir(), "secret-index")
}

func (config *Config) IsComplete() bool {
	return config.DigestedPassphrase != ""
}

func Save(config *Config) error {
	// Marshal the Config struct into YAML
	yamlContent, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write the YAML content to a file
	return os.WriteFile(ConfigPath(), yamlContent, 0644)
}

func LoadConfig(config *Config) error {

	content, err := os.ReadFile(ConfigPath())
	if err != nil {
		return err
	}

	// Unmarshal the YAML content into the Config struct
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}
