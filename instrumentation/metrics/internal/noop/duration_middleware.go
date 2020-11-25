package noop

import (
	"net/http"

	"github.com/puppetlabs/leg/instrumentation/metrics/collectors"
)

type DurationMiddleware struct{}

func (n DurationMiddleware) WithLabels([]collectors.Label) (collectors.DurationMiddleware, error) {
	return n, nil
}

func (n DurationMiddleware) Wrap(next http.Handler) http.Handler { return next }
