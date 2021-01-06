package rand_test

import (
	"sync"
	"testing"

	"github.com/puppetlabs/leg/mathutil/pkg/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMutexGuardedRand(t *testing.T) {
	rng, err := rand.NewTestFactory(rand.NewConstantSeeder(0)).New()
	require.NoError(t, err)

	rng = rand.ThreadSafe(rng)

	// Now we'll read from the test RNG 1,000,000 times across 1,000 Goroutines.
	// If it behaves correctly, the next number it produces will be 1,000,000.
	start := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()

			<-start
			for j := 0; j < 1000; j++ {
				_, err := rand.Uint64(rng)
				require.NoError(t, err)
			}
		}()
	}

	// Tell all Goroutines to start generating.
	close(start)

	// Wait for them to finish.
	wg.Wait()

	rv, err := rand.Uint64(rng)
	require.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), rv)
}
