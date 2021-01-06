package backoff_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstant(t *testing.T) {
	g, err := backoff.Constant(5 * time.Second).New()
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		wait, err := g.Next()
		require.NoError(t, err)
		assert.Equal(t, 5*time.Second, wait, "iteration #%d", i)
	}
}
