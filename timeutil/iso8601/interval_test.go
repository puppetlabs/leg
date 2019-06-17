package iso8601_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/puppetlabs/horsehead/timeutil/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntervalOffset(t *testing.T) {
	ir, err := iso8601.ParseInterval("2016-01-01T00:00:00Z/P1Y")
	require.NoError(t, err)

	e1, err := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	require.NoError(t, err)

	i1 := ir.Offset(3)
	assert.Equal(t, e1, i1.Start())
	assert.Equal(t, ir.Offset(4).Start(), i1.End())

	e2, err := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z")
	require.NoError(t, err)

	i2 := ir.Offset(-2)
	assert.Equal(t, e2, i2.Start())
}

func TestIntervalEquivalence(t *testing.T) {
	ir1, err := iso8601.ParseInterval("2018-10-01T07:00:00Z/P1M")
	require.NoError(t, err)

	ir2, err := iso8601.ParseInterval("P1M/2018-11-01T07:00:00Z")
	require.NoError(t, err)

	ir3, err := iso8601.ParseInterval("2018-10-01T07:00:00Z/2018-11-01T07:00:00Z")
	require.NoError(t, err)

	start, err := time.Parse(time.RFC3339, "2018-10-01T07:00:00Z")
	require.NoError(t, err)

	end, err := time.Parse(time.RFC3339, "2018-11-01T07:00:00Z")
	require.NoError(t, err)

	duration := end.Sub(start)

	for _, ir := range []iso8601.Interval{ir1, ir2, ir3} {
		t.Run(ir.String(), func(t *testing.T) {
			assert.Equal(t, start, ir.Start())
			assert.Equal(t, end, ir.End())
			assert.Equal(t, duration, ir.Duration())
		})
	}
}

func ExampleInterval() {
	e, err := iso8601.ParseInterval("2018-04-01T00:00:00Z/P1M10D")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e.Start().Format(time.RFC3339), e.End().Format(time.RFC3339), e.Duration())

	// Move backward by 2 interval durations.
	e = e.Offset(-2)

	fmt.Println(e.Start().Format(time.RFC3339), e.End().Format(time.RFC3339), e.Duration())

	// Output:
	// 2018-04-01T00:00:00Z 2018-05-11T00:00:00Z 960h0m0s
	// 2018-01-12T00:00:00Z 2018-02-22T00:00:00Z 984h0m0s
}
