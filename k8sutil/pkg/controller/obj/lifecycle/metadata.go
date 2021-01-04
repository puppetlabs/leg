package lifecycle

import (
	"context"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LabelAnnotatableFrom is the type of an entity that can copy labels and
// annotations from the given metadata.
type LabelAnnotatableFrom interface {
	// LabelAnnotateFrom copies the labels from the given object metadata to
	// this entity.
	LabelAnnotateFrom(ctx context.Context, from metav1.Object)
}

// IgnoreNilLabelAnnotatableFrom is an adapter for a label and annotation copier
// that makes sure its delegate has a value before copying.
type IgnoreNilLabelAnnotatableFrom struct {
	LabelAnnotatableFrom
}

// LabelAnnotateFrom copies the labels from the given object metadata to this
// entity.
func (inlaf IgnoreNilLabelAnnotatableFrom) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	if inlaf.LabelAnnotatableFrom == nil || reflect.ValueOf(inlaf.LabelAnnotatableFrom).IsNil() {
		return
	}

	inlaf.LabelAnnotateFrom(ctx, from)
}

var _ LabelAnnotatableFrom = &IgnoreNilLabelAnnotatableFrom{}
