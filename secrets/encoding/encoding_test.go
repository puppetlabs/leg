package encoding

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoding(t *testing.T) {
	var cases = []struct {
		description  string
		value        string
		encodingType encodingType
		expected     string
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

			{
				newResult, err := newED.DecodeSecretValue(value)
				require.NoError(t, err)
				require.Equal(t, c.value, string(newResult))
			}
		})
	}
}
