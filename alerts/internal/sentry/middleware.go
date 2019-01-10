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
		c := m.c.withHTTP(r)
		r = r.WithContext(trackers.NewContextWithCapturer(r.Context(), c))

		defer func() {
			var reporter trackers.Reporter

			rv := recover()
			switch rvt := rv.(type) {
			case nil:
				return
			case error:
				reporter = c.Capture(rvt)
			default:
				reporter = c.CaptureMessage(fmt.Sprint(rvt))
			}

			reporter.Report(r.Context())
		}()

		target.ServeHTTP(w, r)
	})
}
