package sentry

import (
	"fmt"
	"net/http"

	"github.com/puppetlabs/insights-instrumentation/alerts/trackers"
)

type Middleware struct {
	c *Capturer
}

func (m Middleware) WithTags(tags ...trackers.Tag) trackers.Middleware {
	return &Middleware{
		c: m.c.withTags(tags),
	}
}

func (m Middleware) WithUser(u trackers.User) trackers.Middleware {
	return &Middleware{
		c: m.c.withUser(u),
	}
}

func (m Middleware) Wrap(target http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var reporter trackers.Reporter

			rv := recover()
			switch rvt := rv.(type) {
			case nil:
				return
			case error:
				reporter = m.c.withHTTP(r).Capture(rvt)
			default:
				reporter = m.c.withHTTP(r).CaptureMessage(fmt.Sprint(rvt))
			}

			reporter.Report(r.Context())
		}()

		target.ServeHTTP(w, r)
	})
}
