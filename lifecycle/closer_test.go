package lifecycle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/puppetlabs/horsehead/lifecycle"
	"github.com/stretchr/testify/assert"
)

func TestCloserWithoutTrigger(t *testing.T) {
	var i int

	c := lifecycle.NewCloserBuilder().
		Require(func() error {
			i++
			return nil
		}).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-c.Done():
		assert.Fail(t, "closer spuriously terminated")
	case <-time.After(20 * time.Millisecond):
	}

	assert.NoError(t, c.Do(ctx))
	assert.Equal(t, 1, i)
}

func TestCloserInvokesOnTrigger(t *testing.T) {
	ich := make(chan struct{})

	c := lifecycle.NewCloserBuilder().
		When(func(ctx context.Context) error {
			<-ich
			return nil
		}).
		Build()

	select {
	case <-c.Done():
		assert.Fail(t, "closer spuriously terminated")
	default:
	}

	close(ich)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-c.Done():
	case <-ctx.Done():
		assert.Fail(t, "closer did not complete after trigger")
	}

	assert.NoError(t, c.Err())
}

func TestCloserEndsRoutinesOnDo(t *testing.T) {
	var i int

	c := lifecycle.NewCloserBuilder().
		Require(func() error {
			i++
			return nil
		}).
		When(func(ctx context.Context) error {
			<-ctx.Done()
			i++
			return nil
		}).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	assert.NoError(t, c.Do(ctx))
	assert.Equal(t, 2, i)
}

func TestCloserPropagatesErrors(t *testing.T) {
	c := lifecycle.NewCloserBuilder().
		Require(func() error {
			return fmt.Errorf("A")
		}).
		RequireContext(func(ctx context.Context) error {
			return fmt.Errorf("B")
		}).
		When(func(ctx context.Context) error {
			<-ctx.Done()
			return fmt.Errorf("C")
		}).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.Do(ctx)
	assert.NotNil(t, err)
	assert.Len(t, err.Causes(), 3)
}
