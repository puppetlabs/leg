package errmap_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/puppetlabs/leg/errmap/pkg/errmap"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	cause := errors.New("foo")

	err := errmap.MapLast(cause, errmap.MapperFunc(func(err error) error {
		return fmt.Errorf("a: %w", err)
	}))
	err = errmap.MapBefore(err, errmap.MapperFunc(func(err error) error {
		return fmt.Errorf("b: %w", err)
	}))
	err = errmap.MapAfter(err, errmap.MapperFunc(func(err error) error {
		return fmt.Errorf("c: %w", err)
	}))
	err = errmap.MapFirst(err, errmap.MapperFunc(func(err error) error {
		return fmt.Errorf("d: %w", err)
	}))
	err = errmap.MapLast(err, errmap.MapperFunc(func(err error) error {
		return fmt.Errorf("e: %w", err)
	}))

	assert.True(t, errors.Is(err, cause))
	assert.EqualError(t, err, "e: a: c: b: d: foo")
}

type eatingError struct {
	error
}

func (e *eatingError) MapApply(m errmap.Mapper) error {
	mapped := m.Map(e.error)
	return fmt.Errorf("OM NOM NOM: %w", mapped)
}

func TestMapApply(t *testing.T) {
	cause := errors.New("foo")

	err := errmap.Wrap(cause, "a")
	err = errmap.MapLast(err, errmap.MapperFunc(func(err error) error {
		return &eatingError{err}
	}))
	err = errmap.Wrap(err, "b")
	err = errmap.Wrap(err, "c")

	assert.True(t, errors.Is(err, cause))
	assert.EqualError(t, err, "c: OM NOM NOM: b: a: foo")
}

type joinError struct {
	bases []error
}

func (e *joinError) Error() string {
	ps := make([]string, len(e.bases))
	for i, p := range e.bases {
		ps[i] = p.Error()
	}
	return fmt.Sprintf("joined %d error(s): %s", len(e.bases), strings.Join(ps, ", "))
}

func (e *joinError) MapApply(m errmap.Mapper) error {
	ne := &joinError{}
	for _, base := range e.bases {
		next := m.Map(base)

		switch nt := next.(type) {
		case *joinError:
			ne.bases = append(ne.bases, nt.bases...)
		default:
			ne.bases = append(ne.bases, next)
		}
	}
	return ne
}

func TestMapApplyRecursive(t *testing.T) {
	err := errmap.Wrap(errors.New("foo"), "first")
	err = errmap.MapLast(err, errmap.MapperFunc(func(err error) error {
		inner := errmap.MapLast(err, errmap.MapperFunc(func(err error) error {
			return &joinError{bases: []error{fmt.Errorf("inner: %w", err)}}
		}))
		return &joinError{bases: []error{inner, err}}
	}))
	err = errmap.Wrap(err, "outer")

	assert.EqualError(t, err, "joined 2 error(s): outer: inner: first: foo, outer: first: foo")
}
