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

// Label uses any LabelAnnotatableFrom object to assign the given label.
func Label(ctx context.Context, target LabelAnnotatableFrom, key, value string) {
	target.LabelAnnotateFrom(ctx, &metav1.ObjectMeta{
		Labels: map[string]string{key: value},
	})
}

// LabelManagedBy sets the "app.kubernetes.io/managed-by" label to a particular
// value.
func LabelManagedBy(ctx context.Context, target LabelAnnotatableFrom, value string) {
	Label(ctx, target, "app.kubernetes.io/managed-by", value)
}
