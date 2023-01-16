package iso8601_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRecurringInterval(t *testing.T) {
	tests := []struct {
		Value               string
		ExpectedError       bool
		ExpectedRepetitions int
		ExpectedInterval    string
	}{
		{
			Value:               "R26/2015-01-01/PT15M",
			ExpectedRepetitions: 26,
			ExpectedInterval:    "2015-01-01T00:00:00Z/PT15M",
		},
		{
			Value:            "R/P7D/2019-02-28T00:04:15+03:00",
			ExpectedInterval: "P7D/2019-02-28T00:04:15+03:00",
		},
	}
	for _, test := range tests {
		t.Run(test.Value, func(t *testing.T) {
			ri, err := iso8601.ParseRecurringInterval(test.Value)
			if test.ExpectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			n, ok := ri.Repetitions()
			if test.ExpectedRepetitions > 0 {
				require.Equal(t, test.ExpectedRepetitions, n)
				require.Equal(t, true, ok)
			} else {
				require.Equal(t, false, ok)
			}

			require.Equal(t, test.ExpectedInterval, ri.Interval().String())
		})
	}
}

func TestMarshalRecurringInterval(t *testing.T) {
	j := `"R15/PT15M/2018-10-02T00:00:00-08:00"`

	var ri iso8601.RecurringInterval
	require.NoError(t, json.Unmarshal([]byte(j), &ri))

	r, ok := ri.Repetitions()
	assert.True(t, ok)
	assert.Equal(t, 15, r)
	assert.Equal(t, 15*time.Minute, ri.Interval().Duration())

	jn, err := json.Marshal(ri)
	require.NoError(t, err)
	assert.JSONEq(t, j, string(jn))
}

func TestRecurringIntervalCalculatesNextSimple(t *testing.T) {
	iv, err := iso8601.ParseRecurringInterval("R10/2018-10-02T05:00:00Z/PT15M")
	require.NoError(t, err)

	now, err := time.Parse(time.RFC3339, "2018-10-01T23:00:00Z")
	require.NoError(t, err)

	ivn, ok := iv.Next(now)
	assert.True(t, ok)
	assert.Equal(t, iv.Interval().Start(), ivn.Start())

	now = now.Add(8 * time.Hour)

	ivn, ok = iv.Next(now)
	assert.True(t, ok)
	assert.Equal(t, iv.Interval().Start().Add(2*time.Hour), ivn.Start())
	assert.Equal(t, iv.Interval().Offset(8), ivn)

	now = now.Add(25 * time.Minute)

	_, ok = iv.Next(now)
	assert.False(t, ok)
}

func TestRecurringIntervalCalculatesNextYears(t *testing.T) {
	iv, err := iso8601.ParseRecurringInterval("R/2015-10-02T05:00:00Z/P1Y")
	require.NoError(t, err)

	now, err := time.Parse(time.RFC3339, "2015-10-01T23:00:00Z")
	require.NoError(t, err)

	ivn, ok := iv.Next(now)
	assert.True(t, ok)
	assert.Equal(t, iv.Interval().Start(), ivn.Start())

	now = now.AddDate(2, 0, 0)

	ivn, ok = iv.Next(now)
	assert.True(t, ok)
	assert.Equal(t, iv.Interval().Start().AddDate(2, 0, 0), ivn.Start())
	assert.Equal(t, iv.Interval().Offset(2), ivn)
}

func ExampleRecurringInterval() {
	r, err := iso8601.ParseRecurringInterval("R10/2018-10-01T00:00:00Z/P1Y")
	if err != nil {
		log.Fatal(err)
	}

	repetitions, _ := r.Repetitions()
	fmt.Println(repetitions)

	rel1, _ := time.Parse(time.RFC3339, "2019-12-01T00:00:00Z")
	rel2, _ := time.Parse(time.RFC3339, "2029-01-01T00:00:00Z")

	next, _ := r.Next(rel1)
	fmt.Println(next)

	_, ok := r.Next(rel2)
	fmt.Println(ok)

	// Output:
	// 10
	// 2020-10-01T00:00:00Z/P1Y
	// false
}
