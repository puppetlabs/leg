package scheduler

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/puppetlabs/horsehead/v2/netutil"
)

const (
	defaultBackoffMultiplier         = time.Millisecond * 5
	defaultMaxRetries                = 10
	defaultResetRetriesTimerDuration = time.Second * 10
)

// ErrMaxRetriesReached is the error returned by RecoveryDescriptor if the max retries
// have been reached.
var ErrMaxRetriesReached = errors.New("RecoveryDescriptor: max retries have been reached")

// RecoveryDescriptorOptions contains fields that allow backoff and retry parameters
// to be set.
type RecoveryDescriptorOptions struct {
	// BackoffMultiplier is the timing multiplier between attempts using netutil.Backoff.
	BackoffMultiplier time.Duration
	// MaxRetries is the max times the RecoveryDescriptor should attempt to run the delegate
	// descriptor during a reset retries duration.
	MaxRetries int
	// ResetRetriesTimerDuration is the time it takes to reset the retry count when running
	// a delegate descriptor.
	ResetRetriesTimerDuration time.Duration
}

// RecoveryDescriptor wraps a given descriptor so that it restarts if the
// descriptor itself fails. This is useful for descriptors that work off of
// external information (APIs, events, etc.).
type RecoveryDescriptor struct {
	delegate       Descriptor
	backoff        netutil.Backoff
	maxRetries     int
	currentRetries int
	resetDuration  time.Duration
}

var _ Descriptor = &RecoveryDescriptor{}

// runOnce attempts to run the delegate descriptor. It starts a timer that waits for resetRetriesTimerDuration
// that will reset the retry attempt count to 0 if the delegate runs for the duration without an error. This is
// to prevent hours or days from going by, then an error happens again incrementing the count. If this happens 10
// times, the descriptor will seemingly shutdown after 1 retry, causing confusion.
func (rd *RecoveryDescriptor) runOnce(ctx context.Context, pc chan<- Process) (bool, error) {
	timer := time.AfterFunc(rd.resetDuration, func() {
		rd.currentRetries = 0
	})

	err := rd.delegate.Run(ctx, pc)

	// If the timer is already triggered, then this will just return false. So it's
	// fine to call blindly here.
	timer.Stop()

	select {
	case <-ctx.Done():
		return false, err
	default:
	}

	if err != nil {
		log(ctx).Warn("restarting failing descriptor", "descriptor", reflect.TypeOf(rd.delegate).String(), "error", err)
	}

	return true, nil
}

// Run delegates work to another descriptor, catching any errors are restarting
// the descriptor immediately if an error occurs. It never returns an error. It
// only terminates when the context is done.
func (rd *RecoveryDescriptor) Run(ctx context.Context, pc chan<- Process) error {
	for {
		if cont, err := rd.runOnce(ctx, pc); err != nil {
			return err
		} else if !cont {
			break
		}

		if rd.currentRetries == rd.maxRetries {
			log(ctx).Error("max retries reached; stopping descriptor", "descriptor", reflect.TypeOf(rd.delegate).String())
			return ErrMaxRetriesReached
		}

		rd.currentRetries++

		if err := rd.backoff.Backoff(ctx, rd.currentRetries); err != nil {
			return err
		}
	}

	return nil
}

// NewRecoveryDescriptor creates a new recovering descriptor wrapping the given
// delegate descriptor. Default backoff and retry parameters will be used.
func NewRecoveryDescriptor(delegate Descriptor) *RecoveryDescriptor {
	return NewRecoveryDescriptorWithOptions(delegate, RecoveryDescriptorOptions{})
}

// NewRecoveryDescriptorWithOptions creates a new recovering descriptor wrapping the
// given delegate descriptor. It takes RecoveryDescriptorOptions to tune backoff and retry
// parameters.
func NewRecoveryDescriptorWithOptions(delegate Descriptor, opts RecoveryDescriptorOptions) *RecoveryDescriptor {
	if opts.BackoffMultiplier == 0 {
		opts.BackoffMultiplier = defaultBackoffMultiplier
	}

	if opts.MaxRetries == 0 {
		opts.MaxRetries = defaultMaxRetries
	}

	if opts.ResetRetriesTimerDuration == 0 {
		opts.ResetRetriesTimerDuration = defaultResetRetriesTimerDuration
	}

	// TODO migrate to backoff's NextRun once implemented
	backoff := &netutil.ExponentialBackoff{Multiplier: opts.BackoffMultiplier}

	return &RecoveryDescriptor{
		delegate:      delegate,
		backoff:       backoff,
		maxRetries:    opts.MaxRetries,
		resetDuration: opts.ResetRetriesTimerDuration,
	}
}
