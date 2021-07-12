package deduplication

import "fmt"

type bloomFilterBackend interface {
	Set(key string) error
	Exists(key string) (bool, error)
	CheckAndSet(key string) (bool, error)
}

type bloomFilterOptions struct {
	validators []KeyValidatorFunc
}

type bloomFilterOptionsFunc struct {
	f func(o *bloomFilterOptions)
}

func (f bloomFilterOptionsFunc) apply(o *bloomFilterOptions) {
	f.f(o)
}

// BloomFilterOption takes a pointer to bloomFilterOptions and apply's some
// value to one or more fields.
type BloomFilterOption interface {
	apply(*bloomFilterOptions)
}

// WithBloomFilterKeyValidators sets one or more key validators that are used
// to validate keys passed to BloomFilter methods.
func WithBloomFilterKeyValidators(vs ...KeyValidatorFunc) BloomFilterOption {
	return bloomFilterOptionsFunc{
		f: func(o *bloomFilterOptions) {
			o.validators = append(o.validators, vs...)
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
	validators []KeyValidatorFunc
}

// SetKeyAsSeen adds the key to the set of known keys. If you are doing an
// exists check before this method, you probably want to use the atomic
// CheckAndSetKey instead.
func (bf *BloomFilter) SetKeyAsSeen(key string) error {
	if err := bf.validate(key); err != nil {
		return err
	}

	return bf.delegate.Set(key)
}

// KeyHasBeenSeen checks if the key exists in the set of known keys. If you
// want to set they key using SetKeyAsSeen after this, you probably want the
// atomic CheckAndSetKey instead.
func (bf *BloomFilter) KeyHasBeenSeen(key string) (bool, error) {
	return bf.delegate.Exists(key)
}

// CheckAndSetKey checks if the key exists in the set of known keys. If it does
// not exist, then it is added. This method returns whether or not the key
// already existed.
func (bf *BloomFilter) CheckAndSetKey(key string) (bool, error) {
	if err := bf.validate(key); err != nil {
		return false, err
	}

	return bf.delegate.CheckAndSet(key)
}

func (bf *BloomFilter) validate(key string) error {
	for _, fn := range bf.validators {
		if err := fn(key); err != nil {
			return fmt.Errorf("bloom filter: failed to validate key: %w", err)
		}
	}

	return nil
}

// NewBloomFilter takes a bloomFilterBackend and returns a new BloomFilter.
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
