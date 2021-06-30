package deduplication

import "fmt"

type bloomFilterBackend interface {
	set(key string) error
	exists(key string) (bool, error)
}

type bloomFilterOptions struct {
	validators []KeyValidator
}

type bloomFilterOptionsFunc struct {
	f func(o *bloomFilterOptions)
}

func (f bloomFilterOptionsFunc) apply(o *bloomFilterOptions) {
	f.f(o)
}

type BloomFilterOption interface {
	apply(*bloomFilterOptions)
}

func WithBloomFilterKeyValidators(vs ...KeyValidator) BloomFilterOption {
	return bloomFilterOptionsFunc{
		f: func(o *bloomFilterOptions) {
			o.validators = vs
		},
	}
}

// BloomFilter is an implementation of KeySetter and KeyChecker that uses a
// bloom filter data structure to mark keys as seen and to check if keys have
// been seen before. Bloom filters have the possibility for false positives,
// so make sure this is something you can accept before using this
// implementation for data deduplication.
//
// This implementation is geared towards a simple key existance, the key being
// a composite of data points about an object that must be recognized in the
// future.
//
// See https://en.wikipedia.org/wiki/Bloom_filter
type BloomFilter struct {
	delegate   bloomFilterBackend
	validators []KeyValidator
}

func (bf *BloomFilter) SetKeyAsSeen(key string) error {
	for _, v := range bf.validators {
		if err := v.Apply(key); err != nil {
			return fmt.Errorf("bloom filter: failed to validate key: %w", err)
		}
	}

	return bf.delegate.set(key)
}

func (bf *BloomFilter) KeyHasBeenSeen(key string) (bool, error) {
	return bf.delegate.exists(key)
}

func NewBloomFilter(delegate bloomFilterBackend, opts ...BloomFilterOption) *BloomFilter {
	bfo := bloomFilterOptions{}

	for _, o := range opts {
		o.apply(&bfo)
	}

	return &BloomFilter{
		delegate:   delegate,
		validators: bfo.validators,
	}
}
