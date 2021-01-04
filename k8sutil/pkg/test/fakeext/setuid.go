package fakeext

import (
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/testing"
)

type setUIDExtension struct{}

func (setUIDExtension) setUID(obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	accessor.SetUID(types.UID(uuid.New().String()))
	return nil
}

func (sue *setUIDExtension) OnInitial(objs []runtime.Object) error {
	for _, obj := range objs {
		if err := sue.setUID(obj); err != nil {
			return err
		}
	}

	return nil
}

func (sue *setUIDExtension) OnNewFake(f Fake) error {
	f.PrependReactor("create", "*", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
		if ca, ok := action.(testing.CreateActionImpl); ok {
			obj = ca.GetObject()
			err = sue.setUID(obj)
		}
		return
	})
	return nil
}

// SetUIDExtension automatically sets the UID field of the object's metadata to
// a randomly-generated UUID when an object is created.
var SetUIDExtension Extension = &setUIDExtension{}
