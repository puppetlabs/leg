package errmark_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"github.com/stretchr/testify/assert"
)

func TestRules(t *testing.T) {
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	e3 := &net.OpError{}

	m := errmark.NewMarker("test")

	tests := []struct {
		Name     string
		Rule     errmark.Rule
		Error    error
		Expected bool
	}{
		{
			Name:     "Exact match",
			Rule:     errmark.RuleExact(e1),
			Error:    e1,
			Expected: true,
		},
		{
			Name:     "Exact does not match",
			Rule:     errmark.RuleExact(e1),
			Error:    e2,
			Expected: false,
		},
		{
			Name:     "Wrapped exact does not match",
			Rule:     errmark.RuleExact(e1),
			Error:    fmt.Errorf("wrapped: %w", e1),
			Expected: false,
		},
		{
			Name:     "Is match",
			Rule:     errmark.RuleIs(e1),
			Error:    e1,
			Expected: true,
		},
		{
			Name:     "Is does not match",
			Rule:     errmark.RuleIs(e1),
			Error:    e2,
			Expected: false,
		},
		{
			Name:     "Type match",
			Rule:     errmark.RuleType(&net.OpError{}),
			Error:    e3,
			Expected: true,
		},
		{
			Name:     "Type does not match",
			Rule:     errmark.RuleType(&net.OpError{}),
			Error:    e1,
			Expected: false,
		},
		{
			Name:     "Wrapped Is match",
			Rule:     errmark.RuleIs(e1),
			Error:    fmt.Errorf("wrapped: %w", e1),
			Expected: true,
		},
		{
			Name:     "All match",
			Rule:     errmark.RuleAll(errmark.RuleIs(e1), errmark.RuleMarked(m)),
			Error:    errmark.Mark(e1, m),
			Expected: true,
		},
		{
			Name:     "All does not match",
			Rule:     errmark.RuleAll(errmark.RuleIs(e1), errmark.RuleMarked(m)),
			Error:    e1,
			Expected: false,
		},
		{
			Name:     "Always match",
			Rule:     errmark.RuleAlways,
			Error:    e1,
			Expected: true,
		},
		{
			Name:     "Any match",
			Rule:     errmark.RuleAny(errmark.RuleIs(e1), errmark.RuleMarked(m)),
			Error:    errmark.Mark(e2, m),
			Expected: true,
		},
		{
			Name:     "Any does not match",
			Rule:     errmark.RuleAny(errmark.RuleIs(e1), errmark.RuleMarked(m)),
			Error:    e2,
			Expected: false,
		},
		{
			Name:     "Never match",
			Rule:     errmark.RuleNever,
			Error:    e1,
			Expected: false,
		},
		{
			Name:     "Not match",
			Rule:     errmark.RuleNot(errmark.RuleAll(errmark.RuleIs(e1), errmark.RuleMarked(m))),
			Error:    errmark.Mark(e1, m),
			Expected: false,
		},
		{
			Name:     "Not does not match",
			Rule:     errmark.RuleNot(errmark.RuleAll(errmark.RuleIs(e1), errmark.RuleMarked(m))),
			Error:    e1,
			Expected: true,
		},
		{
			Name:     "Predicate match",
			Rule:     errmark.RulePredicate(errmark.RuleAlways, func() bool { return true }),
			Error:    e1,
			Expected: true,
		},
		{
			Name:     "Predicate does not match",
			Rule:     errmark.RulePredicate(errmark.RuleAlways, func() bool { return false }),
			Error:    e1,
			Expected: false,
		},
		{
			Name:     "Marked match",
			Rule:     errmark.RuleMarked(m),
			Error:    errmark.Mark(e1, m),
			Expected: true,
		},
		{
			Name:     "Marked does not match",
			Rule:     errmark.RuleMarked(m),
			Error:    errmark.Mark(e1, errmark.NewMarker("test")),
			Expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Rule.Matches(test.Error))
		})
	}
}
