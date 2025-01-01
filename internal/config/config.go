package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SecretStorePath    string `yaml:"secret_store_path"`
	DigestedPassphrase string `yaml:"digested_passphrase"`
}

func (config *Config) IsComplete() bool {
	return config.SecretStorePath != "" && config.DigestedPassphrase != ""
}

func (config *Config) Save(path string) error {
	// Marshal the Config struct into YAML
	yamlContent, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write the YAML content to a file
	return os.WriteFile(path, yamlContent, 0644)
}

func LoadConfig(path string) (Config, error) {
	// Create an empty Config struct
	var config Config

	content, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	// Unmarshal the YAML content into the Config struct
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
