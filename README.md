# VKE

The VKE library helps in an easy way obtains data from Vault's secrets using Kubernetes auth backend.

# Example of usage

```go
package main

import (
	"fmt"
	"tczekajlo/vke"
)

// Config structure has to contain "Data" structure with fields of our secrets
// All Data structure field has to be "string" type. 
type Config struct {
	Data struct {
		Pool string
	}
}

func main() {
	var config Config

	vke.ReadVaultSecret(&config)

	fmt.Printf("%#v", config)
}
```

Vault data

```
VAULT_ADDR=http://127.0.0.1:8200 vault kv get secret/app-1/config
====== Metadata ======
Key              Value
---              -----
created_time     2019-01-19T09:40:10.150651Z
deletion_time    n/a
destroyed        false
version          1

===== Data =====
Key        Value
---        -----
enabled    true
pool       12
worker     5
```

## Environment variables

The VKE use below environment variables to configuration.

- `APP_ROLE` - application role
- `APP_CONFIG_PATH` - a path to kv secretm, e.g. `secret/data/app-1/config`
- `APP_CONFIG_VERSION` - a version of a key to read, if not given then 1 is a default value