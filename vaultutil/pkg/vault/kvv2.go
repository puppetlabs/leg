package vault

import (
	"context"
	"errors"
	"path"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/puppetlabs/leg/encoding/transfer"
)

var (
	ErrVaultDataNotFound = errors.New("Vault data not found")
)

type KVV2Client struct {
	client     *vaultapi.Client
	enginePath string
	path       string
}

func (c *KVV2Client) In(sub ...string) *KVV2Client {
	return &KVV2Client{
		client:     c.client,
		enginePath: c.enginePath,
		path:       path.Join(c.path, path.Join(sub...)),
	}
}

func (c *KVV2Client) Read(ctx context.Context) (interface{}, error) {
	sec, err := c.client.Logical().Read(c.dataPath())
	if err != nil {
		return nil, err
	} else if sec == nil {
		return nil, ErrVaultDataNotFound
	}

	data, ok := sec.Data["data"].(map[string]interface{})
	if !ok {
		return nil, ErrVaultDataNotFound
	}

	value, found := data["value"]
	if !found {
		return nil, ErrVaultDataNotFound
	}

	return value, nil
}

func (c *KVV2Client) ReadString(ctx context.Context) (string, error) {
	raw, err := c.Read(ctx)
	if err != nil {
		return "", err
	}

	encoded, ok := raw.(string)
	if !ok {
		return "", ErrVaultDataNotFound
	}

	b, err := transfer.DecodeFromTransfer(encoded)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *KVV2Client) Write(ctx context.Context, payload map[string]interface{}) error {
	if _, err := c.client.Logical().Write(c.dataPath(), payload); err != nil {
		return err
	}

	return nil
}

func (c *KVV2Client) Delete(ctx context.Context) error {
	if _, err := c.client.Logical().Delete(c.metadataPath()); err != nil {
		return err
	}

	return nil
}

func (c *KVV2Client) List(ctx context.Context) ([]string, error) {
	ls, err := c.client.Logical().List(c.metadataPath())
	if err != nil {
		return nil, err
	} else if ls == nil {
		return nil, ErrVaultDataNotFound
	}

	ki, ok := ls.Data["keys"].([]interface{})
	if !ok {
		return nil, ErrVaultDataNotFound
	}

	keys := make([]string, len(ki))
	for i, k := range ki {
		keys[i], ok = k.(string)
		if !ok {
			return nil, ErrVaultDataNotFound
		}
	}

	return keys, nil
}

func (c *KVV2Client) dataPath() string {
	return path.Join(c.enginePath, "data", c.path)
}

func (c *KVV2Client) metadataPath() string {
	return path.Join(c.enginePath, "metadata", c.path)
}

func NewKVV2Client(delegate *vaultapi.Client, enginePath string) *KVV2Client {
	return &KVV2Client{
		client:     delegate,
		enginePath: enginePath,
	}
}
