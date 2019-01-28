package noop

import (
	"net/http"

	"github.com/puppetlabs/insights-instrumentation/alerts/trackers"
)

type Middleware struct{}

func (m Middleware) WithTags(tags ...trackers.Tag) trackers.Middleware {
	return m
}

func (m Middleware) WithUser(u trackers.User) trackers.Middleware {
	return m
}

func (m Middleware) Wrap(target http.Handler) http.Handler {
	return target
}
