package errmark_test

import (
	"errors"
	"testing"

	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"github.com/stretchr/testify/assert"
)

func TestShort(t *testing.T) {
	cause := errors.New("foo")

	assert.EqualError(t, errmark.MarkShort(errmark.MarkUser(cause)), "foo")
	assert.EqualError(t, errmark.MarkUser(errmark.MarkShort(cause)), "foo")
}
