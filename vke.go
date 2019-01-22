package vke

import (
	"log"
	"os"
	"tczekajlo/vke/utils"
	"tczekajlo/vke/vault"

	"github.com/mitchellh/mapstructure"
)

const (
	// EnvConfigVersion value is the name of ENV variable which keeps config version to read in Vault
	EnvConfigVersion = "APP_CONFIG_VERSION"
	// EnvConfigRole value is the name of ENV variable which keeps role name require to login to Vault
	EnvConfigRole = "APP_ROLE"
	// EnvConfigPath value is the name of ENV variable which contains a path to read in Vault
	EnvConfigPath = "APP_CONFIG_PATH"
)

// Config structure keeps a basic configuration for Vault path and auth
type Config struct {
	Version string
	Role    string
	Path    string
}

// ReadEnvironment reads configuration information from the environment. If
// there is an error, no configuration value is updated
func (c *Config) ReadEnvironment() error {
	// Parse the environment variables
	if v := os.Getenv(EnvConfigVersion); v != "" {
		c.Version = v
	} else {
		c.Version = "1"
		log.Printf("APP_CONFIG_VERSION is set on default value '1'")
	}

	if v := os.Getenv(EnvConfigPath); v != "" {
		c.Path = v
	} else {
		log.Printf("APP_CONFIG_PATH is empty")
	}

	if v := os.Getenv(EnvConfigRole); v != "" {
		c.Role = v
	} else {
		log.Printf("APP_ROLE is empty")
	}

	return nil
}

// ReadVaultSecret returns data of Vault's kv secret
func ReadVaultSecret(data interface{}) error {
	var config Config

	// Read environment variables
	config.ReadEnvironment()

	vaultClient, err := vault.Client()
	if err != nil {
		log.Fatalf("Error with Vault Client: %s", err)
	}
	kubeAuth := vault.KubernetesAuth{Client: vaultClient}
	kubeAuth.Login(config.Role, utils.ReadServiceAccountToken())

	// Read data from Vault
	secret, err := kubeAuth.Client.Logical().ReadWithData(config.Path, map[string][]string{"version": {config.Version}})

	if secret == nil {
		log.Fatalf("Cannot find data for secret %s", config.Path)
	}

	// Check if a given version of a path is destroyed
	metadata := secret.Data["metadata"].(map[string]interface{})
	if metadata["destroyed"].(bool) {
		log.Fatalf("Version %s of %s secret has been permanently destroyed", config.Version, config.Path)
	}

	// Check if a given version of a path still exists
	if metadata["deletion_time"].(string) != "" {
		log.Fatalf("Cannot find data for path: %s, given version (%s) has been deleted at %s", config.Path, config.Version, metadata["deletion_time"])
	}

	mapstructure.Decode(secret.Data, data)

	return nil
}
