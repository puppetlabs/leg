package transfer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoding(t *testing.T) {
	var cases = []struct {
		description      string
		value            string
		encodingType     EncodingType
		expected         string
		customResultTest func(t *testing.T, encoded string, decoded []byte)
	}{
		{
			description:  "base64 encoding succeeds",
			value:        "super secret token",
			encodingType: Base64EncodingType,
			expected:     "base64:c3VwZXIgc2VjcmV0IHRva2Vu",
		},
		{
			description:  "no encoding succeeds",
			value:        "super secret token",
			encodingType: NoEncodingType,
			expected:     "super secret token",
		},
		{
			description:  "base64 complex values don't loose integrity",
			value:        "super: secret token:12:49:wheel",
			encodingType: Base64EncodingType,
			expected:     "base64:c3VwZXI6IHNlY3JldCB0b2tlbjoxMjo0OTp3aGVlbA==",
		},
		{
			description:  "no encoding complex values don't loose integrity",
			value:        "super: secret token:12:49:wheel",
			encodingType: NoEncodingType,
			expected:     "super: secret token:12:49:wheel",
		},
		{
			description:  "begins with :",
			value:        ":fun time at the park",
			encodingType: NoEncodingType,
			expected:     ":fun time at the park",
		},
		{
			description:  "user encoded base64",
			value:        "c3VwZXIgc2VjcmV0IHRva2Vu",
			encodingType: NoEncodingType,
			expected:     "c3VwZXIgc2VjcmV0IHRva2Vu",
		},
		{
			description: "user encoded base64 wrapped with our base64 encoder",
			// "super secret token" encoded as base64
			value:        "c3VwZXIgc2VjcmV0IHRva2Vu",
			encodingType: Base64EncodingType,
			expected:     "base64:YzNWd1pYSWdjMlZqY21WMElIUnZhMlZ1",
			customResultTest: func(t *testing.T, encoded string, decoded []byte) {
				// tests that a user can encode their own values, have our system wrap it in our own
				// encoding, then when they try to unwrap their encoding, they get the expected value.
				result, err := base64.StdEncoding.DecodeString(string(decoded))
				require.NoError(t, err)

				require.Equal(t, "super secret token", string(result))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			ed := Encoders[c.encodingType]()

			encoded, err := ed.EncodeForTransfer([]byte(c.value))
			require.NoError(t, err)
			require.Equal(t, c.expected, encoded, fmt.Sprintf("encoding result was malformed: %s", encoded))

			typ, value := ParseEncodedValue(encoded)
			require.Equal(t, c.encodingType, typ)

			newED := Encoders[typ]()

			var decoded []byte

			decoded, err = newED.DecodeFromTransfer(value)
			require.NoError(t, err)
			require.Equal(t, c.value, string(decoded))

			if c.customResultTest != nil {
				t.Run("custom result test", func(t *testing.T) {
					c.customResultTest(t, encoded, decoded)
				})
			}
		})
	}
}

func TestHelperFuncs(t *testing.T) {
	var cases = []struct {
		description string
		value       string
		expected    string
	}{
		{
			description: "valid UTF-8 is passed through",
			value:       "super secret token",
			expected:    "super secret token",
		},
		{
			description: "user encoded base64",
			// "super secret token" encoded as base64
			value:    "c3VwZXIgc2VjcmV0IHRva2Vu",
			expected: "c3VwZXIgc2VjcmV0IHRva2Vu",
		},
		{
			description: "invalid UTF-8 is base64-encoded",
			value:       "Hello, \x90\xA2\x8A\x45",
			expected:    "base64:SGVsbG8sIJCiikU=",
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			encoded, err := EncodeForTransfer([]byte(c.value))
			require.NoError(t, err)

			require.Equal(t, c.expected, encoded)

			var decoded []byte

			decoded, err = DecodeFromTransfer(encoded)
			require.NoError(t, err)
			require.Equal(t, c.value, string(decoded))
		})
	}
}

func ExampleEncodeJSON() {
	j, _ := EncodeJSON([]byte("Hello, \x90\xA2\x8A\x45"))
	b, _ := json.Marshal(j)
	fmt.Println(string(b))
	// Output: {"$encoding":"base64","data":"SGVsbG8sIJCiikU="}
}

func ExampleEncodeJSON_plain() {
	j, _ := EncodeJSON([]byte("super secret token"))
	b, _ := json.Marshal(j)
	fmt.Println(string(b))
	// Output: "super secret token"
}
