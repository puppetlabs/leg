package lifecycle

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeleteOptions are the options that a deletable entity is required to support.
type DeleteOptions struct {
	PropagationPolicy client.PropagationPolicy
}

// DeleteOption is a setter for one or more deletion options.
type DeleteOption interface {
	// ApplyToDeleteOptions copies the configuration of this option to the given
	// deletion options.
	ApplyToDeleteOptions(target *DeleteOptions)
}

// ApplyOptions runs each of the given options against this deletion options
// struct.
func (o *DeleteOptions) ApplyOptions(opts []DeleteOption) {
	for _, opt := range opts {
		opt.ApplyToDeleteOptions(o)
	}
}

// DeleteWithPropagationPolicy causes the deletion to use the specified
// propagation policy semantically.
//
// If this option is not used, propagation policy will be determined by the
// entity being deleted.
type DeleteWithPropagationPolicy client.PropagationPolicy

var _ DeleteOption = DeleteWithPropagationPolicy("")

// ApplyToDeleteOptions copies the configuration of this option to the given
// deletion options.
func (dwpp DeleteWithPropagationPolicy) ApplyToDeleteOptions(target *DeleteOptions) {
	target.PropagationPolicy = client.PropagationPolicy(dwpp)
}

// Deleter is the type of an entity that can be deleted from a cluster.
type Deleter interface {
	// Delete removes this entity from the cluster.
	Delete(ctx context.Context, cl client.Client, opts ...DeleteOption) (bool, error)
}
