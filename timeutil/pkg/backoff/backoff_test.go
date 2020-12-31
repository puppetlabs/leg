package backoff_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/mathutil/pkg/rand"
	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var randFactory = rand.NewPCGFactory(rand.OneSeeder)

func TestBuild(t *testing.T) {
	f := backoff.Build(
		backoff.Linear(5*time.Second),
		backoff.MinBound(7*time.Second),
		backoff.MaxBound(30*time.Second),
		backoff.NonSliding,
	)

	b, err := f.New()
	require.NoError(t, err)

	expected := []time.Duration{
		0,
		7 * time.Second,
		10 * time.Second,
		15 * time.Second,
		20 * time.Second,
		25 * time.Second,
		30 * time.Second,
		30 * time.Second,
		30 * time.Second,
	}
	for i, step := range expected {
		wait, err := b.Next()
		require.NoError(t, err)
		assert.Equal(t, step, wait, "step #%d", i)
	}
}
