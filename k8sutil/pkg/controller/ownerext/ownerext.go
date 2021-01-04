package ownerext

import (
	"encoding/json"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type DependencyOf struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Namespace  string    `json:"namespace,omitempty"`
	Name       string    `json:"name"`
	UID        types.UID `json:"uid"`
}

type Manager struct {
	Annotation string
}

func (m *Manager) SetDependencyOf(target metav1.Object, owner lifecycle.TypedObject) error {
	accessor, err := meta.Accessor(owner.Object)
	if err != nil {
		return err
	}

	annotation, err := json.Marshal(DependencyOf{
		APIVersion: owner.GVK.GroupVersion().Identifier(),
		Kind:       owner.GVK.Kind,
		Namespace:  accessor.GetNamespace(),
		Name:       accessor.GetName(),
		UID:        accessor.GetUID(),
	})
	if err != nil {
		return err
	}

	helper.Annotate(target, m.Annotation, string(annotation))
	return nil
}

func (m *Manager) GetDependencyOf(target metav1.Object) (DependencyOf, bool, error) {
	var dep DependencyOf

	annotation := target.GetAnnotations()[m.Annotation]
	if annotation == "" {
		return dep, false, nil
	}

	if err := json.Unmarshal([]byte(annotation), &dep); err != nil {
		return dep, false, err
	}

	return dep, true, nil
}

func (m *Manager) IsDependencyOf(target metav1.Object, owner lifecycle.TypedObject) (bool, error) {
	dep, found, err := m.GetDependencyOf(target)
	if err != nil || !found {
		return found, err
	}

	depGroupVersion, _ := schema.ParseGroupVersion(dep.APIVersion)

	if owner.GVK.Kind != dep.Kind || owner.GVK.GroupVersion() != depGroupVersion {
		return false, nil
	}

	accessor, err := meta.Accessor(owner.Object)
	if err != nil {
		return false, err
	}

	return accessor.GetUID() == dep.UID && accessor.GetNamespace() == dep.Namespace && accessor.GetName() == dep.Name, nil
}

func (m *Manager) NewEnqueueRequestForAnnotatedDependencyOf(ownerType runtime.Object) *EnqueueRequestForAnnotatedDependent {
	return &EnqueueRequestForAnnotatedDependent{
		Manager:   m,
		OwnerType: ownerType,
	}
}

func NewManager(annotation string) *Manager {
	return &Manager{
		Annotation: annotation,
	}
}
