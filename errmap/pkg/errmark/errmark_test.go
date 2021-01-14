package errmark_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/puppetlabs/leg/errmap/pkg/errmap"
	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"github.com/stretchr/testify/assert"
)

func TestMarkSimple(t *testing.T) {
	m := errmark.NewMarker("simple")
	m2 := errmark.NewMarker("also simple")

	err := fmt.Errorf("d'oh")
	assert.False(t, errmark.Marked(err, m))
	assert.Empty(t, errmark.Markers(err).Names())

	err = errmark.Mark(err, m)
	assert.True(t, errmark.Marked(err, m))
	assert.Equal(t, []string{"simple"}, errmark.Markers(err).Names())
	assert.EqualError(t, err, "simple: d'oh")

	err = errmark.Mark(err, m2)
	assert.True(t, errmark.Marked(err, m))
	assert.True(t, errmark.Marked(err, m2))
	assert.Equal(t, []string{"simple", "also simple"}, errmark.Markers(err).Names())
	assert.EqualError(t, err, "simple: also simple: d'oh")

	errU1 := errors.Unwrap(err)
	assert.EqualError(t, errU1, "d'oh")

	errU2 := errors.Unwrap(errmap.Apply(err))
	assert.EqualError(t, errU2, "d'oh")
}

func TestMarkIf(t *testing.T) {
	errCause := errors.New("cause")

	m := errmark.NewMarker("test")

	err := fmt.Errorf("wrapped: %w", errCause)
	err = errmark.MarkIf(err, m, errmark.RuleNever)
	assert.False(t, errmark.Marked(err, m))

	err = errmark.MarkIf(err, m, errmark.RuleIs(errCause))
	assert.True(t, errmark.Marked(err, m))
}

func TestIfUnless(t *testing.T) {
	m := errmark.NewMarker("test")

	tests := []struct {
		Name         string
		Error        error
		Fn           func(err error, rule errmark.Rule, fn func(err error))
		Rule         errmark.Rule
		ExpectCalled bool
	}{
		{
			Name:         "If marked",
			Error:        errmark.Mark(errors.New("d'oh"), m),
			Fn:           errmark.If,
			Rule:         errmark.RuleMarked(m),
			ExpectCalled: true,
		},
		{
			Name:         "If different error",
			Error:        errmark.Mark(errors.New("d'oh"), m),
			Fn:           errmark.If,
			Rule:         errmark.RuleIs(errors.New("nope")),
			ExpectCalled: false,
		},
		{
			Name:         "Unless marked",
			Error:        errmark.Mark(errors.New("d'oh"), m),
			Fn:           errmark.Unless,
			Rule:         errmark.RuleMarked(m),
			ExpectCalled: false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var called bool
			test.Fn(test.Error, test.Rule, func(err error) {
				called = true
			})
			assert.Equal(t, test.ExpectCalled, called)
		})
	}
}
