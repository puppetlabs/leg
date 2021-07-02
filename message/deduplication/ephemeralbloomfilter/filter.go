package ephemeralbloomfilter

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type options struct {
	estimatedItems    uint
	falsePositiveRate float64
}

type optionsFunc struct {
	f func(o *options)
}

func (f optionsFunc) apply(o *options) {
	f.f(o)
}

type Option interface {
	apply(*options)
}

// WithEstimatedItems sets the estimated items expected for the lifetime of the
// bloom filter.
func WithEstimatedItems(items uint) Option {
	return optionsFunc{
		f: func(o *options) {
			o.estimatedItems = items
		},
	}
}

// WithFalsePositiveRate sets the desired false positive rate in the bloom
// filter which lets the implementation decide how many hashing algorithms and
// bit array size to use.
func WithFalsePositiveRate(rate float64) Option {
	return optionsFunc{
		f: func(o *options) {
			o.falsePositiveRate = rate
		},
	}
}

type backend struct {
	bf *bloom.BloomFilter
}

func (e *backend) Set(key string) error {
	e.bf.AddString(key)

	return nil
}

func (e *backend) Exists(key string) (bool, error) {
	return e.bf.TestString(key), nil
}

func (e *backend) CheckAndSet(key string) (bool, error) {
	return e.bf.TestAndAddString(key), nil
}

// New returns a implementation if
// bloomFilterBackend that will not be persisted into storage.
func New(opts ...Option) *backend {
	defaultOpts := options{
		estimatedItems:    1000000,
		falsePositiveRate: 0.01,
	}

	for _, o := range opts {
		o.apply(&defaultOpts)
	}

	return &backend{
		bf: bloom.NewWithEstimates(
			defaultOpts.estimatedItems,
			defaultOpts.falsePositiveRate),
	}
}
