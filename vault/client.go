package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Client return Vault' client
func Client() (*api.Client, error) {
	config := api.DefaultConfig()

	config.ReadEnvironment()

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("err: %s", err)
	}
	return client, nil
}
