package vault

import "errors"

var (
	ErrNoSuchEngineMount = errors.New("engine mount does not exist")
	ErrVaultDataNotFound = errors.New("Vault data not found")
)
