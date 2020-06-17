package stringutil_test

import (
	"testing"

	"github.com/puppetlabs/horsehead/v2/stringutil"
	"github.com/stretchr/testify/assert"
)

func TestUniques(t *testing.T) {
	tt := []struct{ Test, Expected []string }{
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "a", "b", "b", "c", "c"}, []string{"a", "b", "c"}},
		{[]string{"b", "c", "d", "e", "f"}, []string{"b", "c", "d", "e", "f"}},
		{[]string{"w", "x", "y", "y", "y", "y", "y"}, []string{"w", "x", "y"}},
	}
	for _, test := range tt {
		assert.Equal(t, test.Expected, stringutil.Uniques(test.Test), "input: %v", test.Test)
	}
}
