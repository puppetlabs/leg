package encoding

// Encoder encodes a byte slice and returns a string with the encoding type prefixed
type Encoder interface {
	EncodeSecretValue([]byte) (string, error)
}

// Decoder takes a string and decodes it, returning a byte slice or an error
type Decoder interface {
	DecodeSecretValue(string) ([]byte, error)
}

// EncodeDecoder groups Encoder and Decoder to form a type that can both encode and decode
// secret values.
type EncodeDecoder interface {
	Encoder
	Decoder
}
