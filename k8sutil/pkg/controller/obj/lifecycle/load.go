package lifecycle

import (
	"context"
	"reflect"

	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Loader is the type of an entity that can be retrieved from a cluster.
type Loader interface {
	// Load finds this entity in the cluster and populates any necessary fields.
	// If there was an error locating the entity, this function returns false.
	Load(ctx context.Context, cl client.Client) (bool, error)
}

// LoaderFunc allows a function to be used as a loader.
type LoaderFunc func(ctx context.Context, cl client.Client) (bool, error)

var _ Loader = LoaderFunc(nil)

// Load finds this entity in the cluster and populates any necessary fields.
func (lf LoaderFunc) Load(ctx context.Context, cl client.Client) (bool, error) {
	return lf(ctx, cl)
}

// IgnoreNilLoader is an adapter for a loader that makes sure the loader has a
// value before attempting to load from it.
//
// This loader is useful for loading optional dependencies.
type IgnoreNilLoader struct {
	Loader
}

// Load finds this entity in the cluster and populates any necessary fields. If
// the underlying loader is nil or an interface with a nil value, this method
// always returns true with no error.
func (inl IgnoreNilLoader) Load(ctx context.Context, cl client.Client) (bool, error) {
	if inl.Loader == nil || reflect.ValueOf(inl.Loader).IsNil() {
		return true, nil
	}

	return inl.Loader.Load(ctx, cl)
}

// RequiredLoader is an adapter for a loader that ensures the entity exists.
type RequiredLoader struct {
	Loader
}

// Load finds this entity in the cluster and populates any necessary fields. If
// the underlying loader returns false, this loader converts it to an error of
// type *RequiredError.
func (rl RequiredLoader) Load(ctx context.Context, cl client.Client) (bool, error) {
	ok, err := rl.Loader.Load(ctx, cl)
	if err != nil {
		return false, err
	} else if !ok {
		return false, &RequiredError{Loader: rl.Loader}
	}

	return true, nil
}

// Loaders allows multiple loaders to be loaded at once, as if they were a
// single entity.
type Loaders []Loader

var _ Loader = Loaders(nil)

// Load finds this collection of entities in the cluster and loads each one.
//
// If any loader returns false, this method returns false, but it attempts to
// continue loading the remaining entities. If any loader returns an error, the
// error is immediately returned and loading stops.
func (ls Loaders) Load(ctx context.Context, cl client.Client) (bool, error) {
	all := true

	for _, l := range ls {
		if ok, err := l.Load(ctx, cl); err != nil {
			return false, err
		} else if !ok {
			all = false
		}
	}

	return all, nil
}

// RetryLoader is an adapter that continually retries loading until a condition
// is met.
type RetryLoader interface {
	Loader

	// WithWaitOptions sets the options, like a backoff, to use, for waiting on
	// the condition.
	WithWaitOptions(opts ...retry.WaitOption) Loader
}

type retryLoader struct {
	delegate Loader
	mapper   func(ok bool, err error) (bool, error)
}

func (rl *retryLoader) Load(ctx context.Context, cl client.Client) (bool, error) {
	return rl.WithWaitOptions().Load(ctx, cl)
}

func (rl *retryLoader) WithWaitOptions(opts ...retry.WaitOption) Loader {
	return LoaderFunc(func(ctx context.Context, cl client.Client) (ok bool, err error) {
		err = retry.Wait(ctx, func(ctx context.Context) (bool, error) {
			return rl.mapper(rl.delegate.Load(ctx, cl))
		}, opts...)
		return
	})
}

// NewRetryLoader creates a new loader that delegates to the given loader. When
// the Load method is called, its result is mapped using the given condition
// function. This loader only successfully returns when the mapper function
// returns true (even in the case of errors).
func NewRetryLoader(delegate Loader, mapper func(bool, error) (bool, error)) RetryLoader {
	return &retryLoader{
		delegate: delegate,
		mapper:   mapper,
	}
}
