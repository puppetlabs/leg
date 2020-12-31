package retry

import (
	"context"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
)

// WorkFunc is a function that performs an arbitrary operation. If the operation
// needs to be retried for any reason, the function must return false.
type WorkFunc func(ctx context.Context) (bool, error)

// Wait runs a given work function under a context with a specified backoff
// algorithm if the work needs to be retried.
//
// Each time the work is attempted, this function sets its return value to the
// error produced by the work. If the context expires and the work has not
// returned an error, the context error is returned instead.
func Wait(ctx context.Context, bf *backoff.Factory, work WorkFunc) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err == nil {
			err = ctx.Err()
		}
	}()
	defer cancel()

	b, err := bf.New()
	if err != nil {
		return
	}

	bv, err := b.Next()
	if err != nil {
		return
	}

	t := time.NewTimer(bv)
	defer func() { t.Stop() }()

	for {
		select {
		case <-t.C:
		case <-ctx.Done():
			return
		}

		var ok bool
		ok, err = work(ctx)
		if ok {
			return
		}

		bv, err = b.Next()
		if err != nil {
			return
		}

		t.Reset(bv)
	}
}

// WaitAsync runs a given work function in a separate goroutine, but otherwise
// behaves identically to Wait.
func WaitAsync(ctx context.Context, bf *backoff.Factory, work WorkFunc) <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- Wait(ctx, bf, work)
	}()
	return ch
}
