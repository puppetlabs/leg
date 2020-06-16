package stringutil_test

import (
	"testing"

	"github.com/puppetlabs/horsehead/v2/stringutil"
	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	tt := []struct{ Prev, Next, Added, Removed []string }{
		{
			Prev:    []string{"b", "c", "d"},
			Next:    []string{"d", "c", "a"},
			Added:   []string{"a"},
			Removed: []string{"b"},
		},
		{
			Next:  []string{"a", "b", "c", "d"},
			Added: []string{"a", "b", "c", "d"},
		},
		{
			Prev:    []string{"a", "b", "c", "d"},
			Removed: []string{"a", "b", "c", "d"},
		},
		{
			Prev:    []string{"q", "q", "q", "sad"},
			Next:    []string{"happy"},
			Added:   []string{"happy"},
			Removed: []string{"q", "sad"},
		},
		{
			Prev:    []string{"x", "y", "z"},
			Next:    []string{"z", "z"},
			Removed: []string{"x", "y"},
		},
		{
			Prev:    []string{"car", "boat", "train", "plane"},
			Next:    []string{"plane", "automobile", "hot air balloon"},
			Removed: []string{"boat", "car", "train"},
			Added:   []string{"automobile", "hot air balloon"},
		},
	}
	for _, test := range tt {
		added, removed := stringutil.Diff(test.Prev, test.Next)
		assert.Equal(t, test.Added, added)
		assert.Equal(t, test.Removed, removed)
	}
}
