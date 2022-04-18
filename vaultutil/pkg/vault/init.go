package vault

import (
	"context"

	"github.com/google/wire"
	"github.com/puppetlabs/leg/vaultutil/pkg/model"
)

var VaultInitializationManagerProviderSet = wire.NewSet(
	NewVaultInitializationManager,
)

type VaultInitializationManager struct {
	vaultConfig        *VaultConfig
	vaultSystemManager model.VaultSystemManager
}

func (vim *VaultInitializationManager) InitializeVault(
	ctx context.Context, vaultInitializationData *model.VaultInitializationData) error {

	if err := vim.initializeSystem(ctx); err != nil {
		return err
	}

	return vim.initializeData(vaultInitializationData)
}

func (vim *VaultInitializationManager) initializeSystem(ctx context.Context) error {
	credentials := &model.VaultKeys{
		RootToken:  vim.vaultConfig.Token,
		UnsealKeys: []string{vim.vaultConfig.UnsealKey},
	}

	var err error
	if credentials.RootToken == "" {
		credentials, err = vim.vaultSystemManager.GetCredentials(ctx)
		if err != nil {
			return err
		}

		if credentials == nil {
			credentials, err = vim.vaultSystemManager.Initialize(ctx)
			if err != nil {
				return err
			}
		}
	}

	if len(credentials.UnsealKeys) > 0 && credentials.UnsealKeys[0] != "" {
		err = vim.vaultSystemManager.Unseal(credentials)
		if err != nil {
			return err
		}
	}

	vim.vaultSystemManager.SetToken(credentials)

	return nil
}

func (vim *VaultInitializationManager) initializeData(vaultInitializationData *model.VaultInitializationData) error {
	if vaultInitializationData == nil {
		return nil
	}

	if err := vim.vaultSystemManager.RegisterPlugins(vaultInitializationData.Plugins); err != nil {
		return err
	}

	if err := vim.vaultSystemManager.PutPolicies(vaultInitializationData.Policies); err != nil {
		return err
	}

	if err := vim.vaultSystemManager.ConfigureJWTAuthRoles(vaultInitializationData.JWTRoles); err != nil {
		return err
	}

	if err := vim.vaultSystemManager.ConfigureKubernetesAuthRoles(vaultInitializationData.KubernetesRoles); err != nil {
		return err
	}

	if err := vim.vaultSystemManager.EnableSecretEngines(vaultInitializationData.SecretEngines); err != nil {
		return err
	}

	return nil
}

func NewVaultInitializationManager(
	vaultSystemManager model.VaultSystemManager,
	vaultConfig *VaultConfig) *VaultInitializationManager {

	return &VaultInitializationManager{
		vaultConfig:        vaultConfig,
		vaultSystemManager: vaultSystemManager,
	}
}
