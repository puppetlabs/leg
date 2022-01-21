package model

import (
	"context"

	vaultapi "github.com/hashicorp/vault/api"
)

type VaultSystemManager interface {
	ConfigureJWTAuth(ctx context.Context) error
	ConfigureJWTAuthRoles(roles []*VaultJWTRole) error
	ConfigureKubernetesAuth(ctx context.Context) error
	ConfigureKubernetesAuthRoles(roles []*VaultKubernetesRole) error
	CreateCredentials(ctx context.Context, vaultKeys *VaultKeys) error
	CreateTransitKey(mount, key string) error
	EnableJWTAuth() error
	EnableKubernetesAuth() error
	EnableSecretEngines(secretEngines []*VaultSecretEngine) error
	GetAuthMethod(name string) (*vaultapi.AuthMount, error)
	GetCredentials(ctx context.Context) (*VaultKeys, error)
	GetKubernetesAuthConfig(ctx context.Context) (*KubernetesCertificateData, error)
	Initialize(ctx context.Context) (*VaultKeys, error)
	IsJWTAuthEnabled() (bool, error)
	IsKubernetesAuthEnabled() (bool, error)
	PutPolicies(policies []*VaultPolicy) error
	RegisterPlugins(plugins []*vaultapi.RegisterPluginInput) error
	SetToken(vaultKeys *VaultKeys)
	Unseal(vaultKeys *VaultKeys) error
}
