//go:build wireinject
// +build wireinject

package test

import (
	"context"

	"github.com/google/wire"
	vaultapi "github.com/hashicorp/vault/api"
	fake "github.com/puppetlabs/leg/k8sutil/pkg/test/controller/fake"
	"github.com/puppetlabs/leg/vaultutil/pkg/model"
	"github.com/puppetlabs/leg/vaultutil/pkg/vault"
)

func NewTestVaultSystemManager(ctx context.Context, cfg *vault.VaultConfig, vaultClient *vaultapi.Client) (*vault.VaultSystemManager, func(), error) {
	panic(wire.Build(
		fake.NewClient,
		vault.VaultSystemManagerProviderSet,
	))
}

func NewTestVaultInitializationManager(ctx context.Context, cfg *vault.VaultConfig, vaultClient *vaultapi.Client) (*vault.VaultInitializationManager, func(), error) {
	panic(wire.Build(
		fake.NewClient,
		vault.VaultSystemManagerProviderSet,
		vault.VaultInitializationManagerProviderSet,
		wire.Bind(new(model.VaultSystemManager), new(*vault.VaultSystemManager)),
	))
}
