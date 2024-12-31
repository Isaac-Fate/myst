package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Isaac-Fate/myst/internal/config"
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	config := config.Config{
		PasswordStorePath:  "/path/to/password-store.db",
		DigestedPassphrase: "xxxx",
	}

	yamlContent, err := yaml.Marshal(config)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("yaml content:\n%s\n", string(yamlContent))
}
func TestWriteConfig(t *testing.T) {
	config := config.Config{
		PasswordStorePath:  "/path/to/password-store.db",
		DigestedPassphrase: "xxxx",
	}

	yamlContent, err := yaml.Marshal(&config)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("yaml content:\n%s\n", string(yamlContent))

	err = os.WriteFile("myst.yaml", yamlContent, 0644)
	if err != nil {
		t.Error(err)
	}
}
