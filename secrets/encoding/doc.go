/*
Package encoding provides an interface for encoding and decoding secret
values for storage. The utility in this package is transparent to the user
and it is used to maintain byte integrity on secret values used in workflows.

Basic use when encoding a value:
	encoder := encoding.Encoders[encoding.DefaultEncodingType]()

	result, err := encoder.EncodeSecretValue([]byte("super secret token"))
	if err != nil {
		// handle error
	}

Basic use when decoding a value:
	encodingType, value := encoding.ParseEncodedValue("base64:c3VwZXIgc2VjcmV0IHRva2Vu")
	encoder := encoding.Encoders[encoderType]()

	result, err := encoder.DecodeSecretValue(value)
	if err != nil {
		// handle error
	}
*/
package encoding
