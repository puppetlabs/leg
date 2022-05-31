package evaluate

import "fmt"

var (
	ArrayAnnotation     = Annotation{Name: "leg.relay.sh/array", Description: "array"}
	ExpansionAnnotation = Annotation{Name: "leg.relay.sh/expansion", Description: "expansion"}
	ObjectAnnotation    = Annotation{Name: "leg.relay.sh/object", Description: "object"}
)

func ArrayIndexAnnotation(idx int) Annotation {
	return Annotation{
		Name:        "leg.relay.sh/array.index",
		Attributes:  map[string]any{"index": idx},
		Description: fmt.Sprintf("index %d", idx),
	}
}

func ObjectKeyAnnotation(key string) Annotation {
	return Annotation{
		Name:        "leg.relay.sh/object.key",
		Attributes:  map[string]any{"key": key},
		Description: fmt.Sprintf("key %s", key),
	}
}
