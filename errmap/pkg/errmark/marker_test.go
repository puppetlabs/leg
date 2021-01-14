package errmark_test

import (
	"testing"

	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"github.com/stretchr/testify/assert"
)

func TestMarkerSet(t *testing.T) {
	m1 := errmark.NewMarker("m1")
	m2 := errmark.NewMarker("m2")

	ms1 := errmark.NewMarkerSet(m1)
	ms2 := errmark.NewMarkerSet(m2)

	assert.True(t, ms1.Has(m1))
	assert.False(t, ms1.Has(m2))
	assert.Equal(t, []string{"m1"}, ms1.Names())
	assert.Equal(t, []string{"m2"}, ms2.Names())

	ms3 := ms1.Merge(ms2)
	assert.True(t, ms3.Has(m1))
	assert.True(t, ms3.Has(m2))
	assert.Equal(t, []string{"m1", "m2"}, ms3.Names())
}
