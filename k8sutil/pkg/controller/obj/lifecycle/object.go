package lifecycle

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TypedObject is a Kubernetes runtime object with its resource
// group-version-kind attached from its schema.
type TypedObject struct {
	GVK    schema.GroupVersionKind
	Object runtime.Object
}
