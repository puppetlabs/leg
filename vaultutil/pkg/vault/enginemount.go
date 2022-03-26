package vault

import (
	"path"

	"github.com/hashicorp/vault/api"
)

func CheckNormalizeEngineMount(client *api.Client, mount string) (string, error) {
	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return "", err
	}

	normalizedMount := path.Clean(mount) + `/`

	if _, ok := mounts[normalizedMount]; !ok {
		return "", ErrNoSuchEngineMount
	}

	return normalizedMount, nil
}
