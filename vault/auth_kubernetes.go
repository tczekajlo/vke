package vault

import (
	"log"
	"tczekajlo/vke/utils"

	"github.com/hashicorp/vault/api"
)

// KubernetesAuth structure keeps Vault's client
type KubernetesAuth struct {
	Client *api.Client
}

// Login obtains client token using Vault Kubernetes auth backend
func (c *KubernetesAuth) Login(role, jwt string) error {
	path := "auth/kubernetes/login"
	secret, err := c.Client.Logical().Write(path, map[string]interface{}{
		"role": role,
		"jwt":  utils.ReadServiceAccountToken(),
	})
	if err != nil {
		log.Fatalf("Error writing data to %s: %s", path, err)
		return err
	}

	c.Client.SetToken(secret.Auth.ClientToken)
	return nil
}
