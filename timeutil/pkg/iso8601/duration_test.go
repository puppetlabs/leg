package iso8601_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDurationBeforeAfter(t *testing.T) {
	rel, err := iso8601.ParseDuration("P1Y")
	require.NoError(t, err)

	// Leap year.
	y1, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	require.NoError(t, err)

	// Not leap year.
	y2, err := time.Parse(time.RFC3339, "2017-01-01T00:00:00Z")
	require.NoError(t, err)

	assert.NotEqual(t, rel.After(y1), rel.After(y2))
	assert.Equal(t, -rel.After(y1), rel.Before(y2))
}

func ExampleDuration() {
	rel1, _ := time.Parse(time.RFC3339, "2018-02-01T00:00:00Z")
	rel2, _ := time.Parse(time.RFC3339, "2018-03-01T00:00:00Z")

	d, err := iso8601.ParseDuration("P1MT6H1.5M")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(d.After(rel1))
	fmt.Println(d.After(rel2))

	// Output:
	// 678h1m30s
	// 750h1m30s
}
