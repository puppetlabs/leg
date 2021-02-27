package manifest

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PatcherFunc is the type of an object patcher than can be applied by Parse.
type PatcherFunc func(obj Object, gvk *schema.GroupVersionKind)

// FixupPatcher fixes common issues with YAML files that are usually remediated
// by `kubectl apply` automatically.
func FixupPatcher(obj Object, gvk *schema.GroupVersionKind) {
	switch t := obj.(type) {
	case *appsv1.Deployment:
		// SSA has marked "protocol" is required but basically everyone expects
		// it to default to TCP.
		for i := range t.Spec.Template.Spec.Containers {
			container := &t.Spec.Template.Spec.Containers[i]

			for j, port := range container.Ports {
				if port.Protocol != "" {
					continue
				}

				container.Ports[j].Protocol = corev1.ProtocolTCP
			}
		}
	case *corev1.Service:
		// Same for services.
		for i, port := range t.Spec.Ports {
			if port.Protocol != "" {
				continue
			}

			t.Spec.Ports[i].Protocol = corev1.ProtocolTCP
		}
	}
}

// DefaultNamespacePatcher sets the namespace metadata field of namespace-scoped
// objects to the given namespace name if the field is not specified.
func DefaultNamespacePatcher(mapper meta.RESTMapper, namespace string) PatcherFunc {
	return func(obj Object, gvk *schema.GroupVersionKind) {
		// Namespace already set?
		if obj.GetNamespace() != "" {
			return
		}

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return
		}

		// Does this resource even take a namespace?
		if mapping.Scope.Name() != meta.RESTScopeNameNamespace {
			return
		}

		obj.SetNamespace(namespace)
	}
}
