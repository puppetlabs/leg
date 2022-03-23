//go:build wireinject
// +build wireinject

package vault

import (
	"context"

	"github.com/google/wire"
	vaultapi "github.com/hashicorp/vault/api"
	fake "github.com/puppetlabs/leg/vaultutil/pkg/fake"
	"github.com/puppetlabs/leg/vaultutil/pkg/model"
)

func NewTestVaultSystemManager(ctx context.Context, cfg *VaultConfig, vaultClient *vaultapi.Client) (*VaultSystemManager, func(), error) {
	panic(wire.Build(
		fake.FakeKubeProviderSet,
		VaultSystemManagerProviderSet,
	))
}

func NewTestVaultInitializationManager(ctx context.Context, cfg *VaultConfig, vaultClient *vaultapi.Client) (*VaultInitializationManager, func(), error) {
	panic(wire.Build(
		fake.FakeKubeProviderSet,
		VaultSystemManagerProviderSet,
		VaultInitializationManagerProviderSet,
		wire.Bind(new(model.VaultSystemManager), new(*VaultSystemManager)),
	))
}
