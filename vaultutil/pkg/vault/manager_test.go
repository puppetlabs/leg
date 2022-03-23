package vault_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/puppetlabs/leg/vaultutil/pkg/model"
	vaultutil "github.com/puppetlabs/leg/vaultutil/pkg/vault"
	"github.com/stretchr/testify/require"
)

func TestCoreInitialization(t *testing.T) {

	tcs := []struct {
		name string
		init *model.VaultInitializationData
	}{
		{
			name: "core initialization",
			init: &model.VaultInitializationData{
				SecretEngines: []*model.VaultSecretEngine{
					{
						Path:        uuid.NewString(),
						Type:        model.VaultSecretEngineTypeKVV2.String(),
						Description: uuid.NewString(),
					},
					{
						Path:        uuid.NewString(),
						Type:        model.VaultSecretEngineTypeKVV2.String(),
						Description: uuid.NewString(),
					},
				},
			},
		},
	}

	cfg := &vaultutil.VaultConfig{
		Name:      uuid.NewString(),
		Namespace: uuid.NewString(),
	}

	ctx := context.Background()

	for _, test := range tcs {
		t.Run(test.name, func(t *testing.T) {
			err := WithVaultCore(t, ctx, func(vaultCore *vault.Core) {
				err := WithVaultServer(t, vaultCore, func(vaultClient *vaultapi.Client) {
					err := WithVaultInitializationManager(ctx, cfg, vaultClient, func(vi *vaultutil.VaultInitializationManager) {
						for i := 0; i < 3; i++ {
							err := vi.InitializeVault(ctx, test.init)
							require.NoError(t, err)
						}
					})
					require.NoError(t, err)
				})
				require.NoError(t, err)
			})
			require.NoError(t, err)
		})
	}
}

func WithVaultInitializationManager(ctx context.Context, cfg *vaultutil.VaultConfig, vaultClient *vaultapi.Client, fn func(vsm *vaultutil.VaultInitializationManager)) error {
	vi, cleanup, err := vaultutil.NewTestVaultInitializationManager(ctx, cfg, vaultClient)
	if err != nil {
		return err
	}

	defer cleanup()

	fn(vi)

	return nil
}

func WithVaultSystemManager(ctx context.Context, cfg *vaultutil.VaultConfig, vaultClient *vaultapi.Client, fn func(vsm *vaultutil.VaultSystemManager)) error {
	vsm, cleanup, err := vaultutil.NewTestVaultSystemManager(ctx, cfg, vaultClient)
	if err != nil {
		return err
	}

	defer cleanup()

	fn(vsm)

	return nil
}

func WithVaultCore(t *testing.T, ctx context.Context, fn func(core *vault.Core)) error {
	core := NewTestVaultCore(t)
	fn(core)

	return nil
}

func WithVaultServer(t *testing.T, core *vault.Core, fn func(vc *vaultapi.Client)) error {
	vc, cleanup, err := NewTestVaultServer(t, core)
	if err != nil {
		return err
	}

	defer cleanup()

	fn(vc)

	return nil
}

func NewTestVaultCore(t *testing.T) *vault.Core {
	return vault.TestCore(t)
}

func NewTestVaultServer(t *testing.T, core *vault.Core) (*vaultapi.Client, func(), error) {
	ln, addr := http.TestServer(t, core)

	conf := vaultapi.DefaultConfig()
	conf.Address = addr

	c, err := vaultapi.NewClient(conf)
	if err != nil {
		return nil, nil, err
	}

	return c, func() {
		ln.Close()
	}, nil
}
