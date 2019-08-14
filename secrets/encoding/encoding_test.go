package encoding

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoding(t *testing.T) {
	var cases = []struct {
		description      string
		value            string
		encodingType     encodingType
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
				t.Run("custom result test", func(t *testing.T) {
					// tests that a user can encode their own values, have our system wrap it in our own
					// encoding, then when they try to unwrap their encoding, they get the expected value.
					result, err := base64.StdEncoding.DecodeString(string(decoded))
					require.NoError(t, err)

					require.Equal(t, "super secret token", string(result))
				})
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			ed := Encoders[c.encodingType]()

			result, err := ed.EncodeSecretValue([]byte(c.value))
			require.NoError(t, err)
			require.Equal(t, c.expected, result, fmt.Sprintf("result was malformed: %s", result))

			typ, value := ParseEncodedValue(result)
			require.Equal(t, c.encodingType, typ)

			newED := Encoders[typ]()

			var newResult []byte

			newResult, err = newED.DecodeSecretValue(value)
			require.NoError(t, err)
			require.Equal(t, c.value, string(newResult))

			if c.customResultTest != nil {
				c.customResultTest(t, result, newResult)
			}
		})
	}
}
