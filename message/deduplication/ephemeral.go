package deduplication

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type ephemeralBloomFilterBackend struct {
	bf *bloom.BloomFilter
}

func (e *ephemeralBloomFilterBackend) set(key string) error {
	e.bf.AddString(key)

	return nil
}

func (e *ephemeralBloomFilterBackend) exists(key string) (bool, error) {
	return e.bf.TestString(key), nil
}

func NewEphemeralBloomFilterBackend() *ephemeralBloomFilterBackend {
	return &ephemeralBloomFilterBackend{
		bf: bloom.NewWithEstimates(1000000, 0.01),
	}
}
