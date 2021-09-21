package manifest

import (
	"bytes"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Parse reads a stream of YAML documents from the given reader and parses them
// into Kubernetes objects according to the given runtime scheme.
//
// For each object loaded, the specified list of patchers is run. FixupPatcher
// is automatically run and does not need to be present in the list of patchers.
func Parse(scheme *runtime.Scheme, r io.Reader, patchers ...PatcherFunc) ([]Object, error) {
	patchers = append(patchers, FixupPatcher)

	d := yaml.NewYAMLOrJSONDecoder(r, 4096)

	// This lets us convert input documents.
	deserializer := serializer.NewCodecFactory(scheme).UniversalDeserializer()

	// The objects to create.
	var objs []Object

	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}

		robj, gvk, err := deserializer.Decode(ext.Raw, nil, nil)
		if err != nil {
			return nil, err
		}

		obj, ok := robj.(Object)
		if !ok {
			return nil, fmt.Errorf("object of type %T is missing metadata", obj)
		}

		for _, patcher := range patchers {
			patcher(obj, gvk)
		}

		objs = append(objs, obj)
	}

	return objs, nil
}
