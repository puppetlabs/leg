package manifest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

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

	decoder := yaml.NewDocumentDecoder(ioutil.NopCloser(r))
	defer decoder.Close()

	// Copy buffer; we can't use io.Copy because of the weird semantics of the
	// document decoder in how it returns ErrShortBuffer.
	buf := make([]byte, 32*1024)

	// This lets us convert input documents.
	deserializer := serializer.NewCodecFactory(scheme).UniversalDeserializer()

	// The objects to create.
	var objs []Object

	var stop bool
	for !stop {
		var doc bytes.Buffer

		for {
			nr, err := decoder.Read(buf)
			if nr > 0 {
				if nw, err := doc.Write(buf[:nr]); err != nil {
					return nil, err
				} else if nw != nr {
					return nil, io.ErrShortWrite
				}
			}

			if err == io.ErrShortBuffer {
				// More document to read, keep going.
			} else if err == io.EOF {
				// End of the entire stream.
				stop = true
				break
			} else if err != nil {
				return nil, err
			} else {
				// End of this loop, but we have another document ahead.
				break
			}
		}

		b := doc.Bytes()
		if len(bytes.TrimSpace(b)) == 0 {
			// Empty document.
			continue
		}

		robj, gvk, err := deserializer.Decode(b, nil, nil)
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
