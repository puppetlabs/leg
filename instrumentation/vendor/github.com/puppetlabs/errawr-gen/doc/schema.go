package doc

import "github.com/xeipuuv/gojsonschema"

var (
	Schema, _ = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(MustAsset("schemas/v1/errors.json")))
)
