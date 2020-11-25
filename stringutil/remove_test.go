package stringutil_test

import (
	"testing"

	"github.com/puppetlabs/leg/stringutil"
	"github.com/stretchr/testify/assert"
)

func TestRemoveAll(t *testing.T) {
	tt := []struct{ S, CutSet, Expected string }{
		{
			S:        `this. is ^a "test"`,
			CutSet:   `.,:;/\'"=*!?#$&+^|~<>(){}[]` + "`",
			Expected: "this is a test",
		},
		{
			S:        "test",
			CutSet:   "@#$%^",
			Expected: "test",
		},
	}
	for _, test := range tt {
		assert.Equal(t, test.Expected, stringutil.RemoveAll(test.S, test.CutSet), "for string %q with cutset %q", test.S, test.CutSet)
	}
}
