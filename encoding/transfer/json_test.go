package transfer_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/stretchr/testify/require"
)

type testJSONHexEncoding struct{}

func (testJSONHexEncoding) EncodeForTransfer(value []byte) (string, error) {
	return fmt.Sprintf("hex:%s", hex.EncodeToString(value)), nil
}

func (testJSONHexEncoding) EncodeJSON(value []byte) (transfer.JSONOrStr, error) {
	return transfer.JSONOrStr{JSON: transfer.JSON{
		EncodingType: transfer.EncodingType("hex"),
		Data:         hex.EncodeToString(value),
	}}, nil
}

func (testJSONHexEncoding) DecodeFromTransfer(value string) ([]byte, error) {
	return hex.DecodeString(value)
}

var testJSONFactories = transfer.EncodeDecoderFactories{
	transfer.EncodingType("hex"): func() transfer.EncodeDecoder { return &testJSONHexEncoding{} },
}

func TestJSONUnmarshal(t *testing.T) {
	var cases = []struct {
		description string
		json        string
		expected    string
		factories   transfer.EncodeDecoderFactories
		err         error
	}{
		{
			description: "Base64 encoding succeeds",
			json:        `{"$encoding": "base64", "data": "c3VwZXIgc2VjcmV0IHRva2Vu"}`,
			expected:    "super secret token",
		},
		{
			description: "Explicit empty encoding succeeds",
			json:        `{"$encoding": "", "data": "blah blah blee bloo"}`,
			expected:    "blah blah blee bloo",
		},
		{
			description: "Invalid encoding errors",
			json:        `{"$encoding": "invalid", "data": "blah blah blee bloo"}`,
			err:         transfer.ErrUnknownEncodingType,
		},
		{
			description: "Custom encoder factory",
			json:        `{"$encoding": "hex", "data": "48656c6c6f20476f7068657221"}`,
			expected:    "Hello Gopher!",
			factories:   testJSONFactories,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			j := transfer.JSON{
				Factories: c.factories,
			}

			require.NoError(t, json.Unmarshal([]byte(c.json), &j))
			b, err := j.Decode()
			if c.err != nil {
				require.Equal(t, c.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, c.expected, string(b))
		})
	}
}

func TestJSONOrStrMarshalUnmarshal(t *testing.T) {
	var cases = []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "Properly encodes utf8 strings",
			input:       "This is a normal string",
			expected:    `"This is a normal string"`,
		},
		{
			description: "Properly encodes non-utf8 strings",
			input:       "Hello, \x90\xA2\x8A\x45",
			expected:    `{"$encoding":"base64","data":"SGVsbG8sIJCiikU="}`,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			j, err := transfer.EncodeJSON([]byte(c.input))
			require.NoError(t, err)

			js, err := json.Marshal(j)
			require.NoError(t, err)
			require.JSONEq(t, c.expected, string(js))

			var ju transfer.JSONOrStr
			require.NoError(t, json.Unmarshal(js, &ju))

			d, err := ju.Decode()
			require.Equal(t, c.input, string(d))
		})
	}
}
