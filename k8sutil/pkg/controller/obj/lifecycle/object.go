package lifecycle

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TypedObject is a Kubernetes runtime object with its resource
// group-version-kind attached from its schema.
type TypedObject struct {
	GVK    schema.GroupVersionKind
	Object client.Object
}
