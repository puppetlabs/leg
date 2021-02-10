package types_test

import (
	"encoding/json"
	"testing"

	"github.com/puppetlabs/leg/jsonutil/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestKeyPathSetByJSONArray(t *testing.T) {
	template := `{"key_path": ["/tmp/1", "/tmp/2"]}`

	var config struct {
		KeyPath types.StrOrStrSlice `json:"key_path"`
	}
	assert.NoError(t, json.Unmarshal([]byte(template), &config))
	assert.Len(t, config.KeyPath, 2)
	assert.Equal(t, "/tmp/1", config.KeyPath[0])
	assert.Equal(t, "/tmp/2", config.KeyPath[1])
}

func TestKeyPathSetByJSONString(t *testing.T) {
	template := `{"key_path": "/tmp/1"}`

	var config struct {
		KeyPath types.StrOrStrSlice `json:"key_path"`
	}
	assert.NoError(t, json.Unmarshal([]byte(template), &config))
	assert.Len(t, config.KeyPath, 1)
	assert.Equal(t, "/tmp/1", config.KeyPath[0])
}
