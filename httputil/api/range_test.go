package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newInt64(x int64) *int64 {
	return &x
}

func TestScanRangeHeader(t *testing.T) {
	tests := []struct {
		Header   string
		Error    error
		Expected *RangeHeader
	}{
		{
			Header: ` bytes = 1 - 3 , 10 - `,
			Expected: &RangeHeader {
				Unit: "bytes",
				Specs: []RangeSpec {
					{
						First: newInt64(1),
						Last:  newInt64(3),
					},
					{
						First: newInt64(10),
					},
				},
			},
		},
		{
			Header: ` bytes = - 1000`,
			Expected: &RangeHeader {
				Unit: "bytes",
				Specs: []RangeSpec {
					{
						SuffixLength: newInt64(1000),
					},
				},
			},
		},
		{
			Header: ` lines =1-2`,
			Error:  fmt.Errorf(`Unsupported Range header unit=lines`),
		},
		{
			Header: `bytes=-`,
			Error:  fmt.Errorf(`Invalid Range header, expected more than just a '-'`),
		},
	}
	for _, test := range tests {
		spec, err := ScanRangeHeader(test.Header)

		assert.Equal(t, test.Error, err, "for header %q", test.Header)
		assert.Equal(t, test.Expected, spec, "for header %q", test.Header)
	}
}
