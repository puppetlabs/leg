package redisbloomfilter

import (
	"fmt"

	redisbloom "github.com/RedisBloom/redisbloom-go"
	"github.com/gomodule/redigo/redis"
)

const (
	defaultCapacity   = uint64(1_000_000)
	defaultErrorRate  = float64(0.1)
	defaultClientName = "puppetlabs-leg"
)

type FilterName string

type options struct {
	clientName string
	capacity   uint64
	errorRate  float64
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

func WithClientName(name string) Option {
	return optionsFunc{
		f: func(o *options) {
			o.clientName = name
		},
	}
}

func WithCapacity(capacity uint64) Option {
	return optionsFunc{
		f: func(o *options) {
			o.capacity = capacity
		},
	}
}

func WithErrorRate(rate float64) Option {
	return optionsFunc{
		f: func(o *options) {
			o.errorRate = rate
		},
	}
}

type backend struct {
	c          *redisbloom.Client
	filterName string
	capacity   uint64
	errorRate  float64
}

// Set adds the key to the bloom filter under the configured filter name.
func (b *backend) Set(key string) error {
	_, err := b.c.Add(b.filterName, key)
	return err
}

// Exists checks if the key exists under the configured filter name.
func (b *backend) Exists(key string) (bool, error) {
	return b.c.Exists(b.filterName, key)
}

// Check and set checks if the key exists under the configured filter name and
// adds it if it isn't. This method returns whether or not the key already
// existed.
func (b *backend) CheckAndSet(key string) (bool, error) {
	inserted, err := b.c.Add(b.filterName, key)

	return !inserted, err
}

// New takes a filter name, redis connection pool, some options and returns a
// configured bloomFilterBackend or an error. The redis server this client
// connects to must have the RedisBloom module installed. A filter name is
// required and must not be empty.
func New(name FilterName, pool *redis.Pool, opts ...Option) (*backend, error) {
	defaultOpts := options{
		clientName: defaultClientName,
		capacity:   defaultCapacity,
		errorRate:  defaultErrorRate,
	}

	for _, o := range opts {
		o.apply(&defaultOpts)
	}

	filterName := string(name)

	if filterName == "" {
		return nil, fmt.Errorf("filter name is required and cannot be empty")
	}

	client := redisbloom.NewClientFromPool(pool, defaultOpts.clientName)

	// reserve the filter if it hasn't been reserved yet
	err := client.Reserve(filterName, defaultOpts.errorRate, defaultOpts.capacity)
	if err != nil {
		if rerr, ok := err.(redis.Error); !ok || rerr.Error() != "ERR item exists" {
			return nil, fmt.Errorf("failed to reserve bloom filter %s: %w", filterName, err)
		}
	}

	return &backend{
		c:          client,
		filterName: filterName,
		capacity:   defaultOpts.capacity,
		errorRate:  defaultOpts.errorRate,
	}, nil
}
