// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package vault

import (
	"context"
	"github.com/hashicorp/vault/api"
	"github.com/puppetlabs/leg/vaultutil/pkg/fake"
)

// Injectors from wire.go:

func NewTestVaultSystemManager(ctx context.Context, cfg *VaultConfig, vaultClient *api.Client) (*VaultSystemManager, func(), error) {
	client := fake.NewFakeKubeClient()
	vaultSystemManager := NewVaultSystemManager(vaultClient, client, cfg)
	return vaultSystemManager, func() {
	}, nil
}

func NewTestVaultInitializationManager(ctx context.Context, cfg *VaultConfig, vaultClient *api.Client) (*VaultInitializationManager, func(), error) {
	client := fake.NewFakeKubeClient()
	vaultSystemManager := NewVaultSystemManager(vaultClient, client, cfg)
	vaultInitializationManager := NewVaultInitializationManager(vaultSystemManager, cfg)
	return vaultInitializationManager, func() {
	}, nil
}