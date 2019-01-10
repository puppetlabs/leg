package trackers

import (
	"context"
)

type Capturer interface {
	// WithNewTrace causes any captures to automatically get a stack trace
	// associated to them.
	WithNewTrace() Capturer

	// WithAppPackages causes the listed packages (and any child packages) to be
	// highlighted as part of the current application in stack traces.
	WithAppPackages(packages []string) Capturer

	// WithUser adds the given user information to any errors reported.
	WithUser(u User) Capturer

	// WithTags adds the given tags to any errors reported.
	WithTags(tags ...Tag) Capturer

	// Try runs the given function, and if a panic occurs, captures and reports
	// it. It returns the recovered value of the panic, or nil if no panic
	// occurred.
	Try(ctx context.Context, fn func()) interface{}

	// Capture captures the given error for reporting.
	Capture(err error) Reporter

	// CaptureMessage converts the given message to an error and captures it.
	CaptureMessage(message string) Reporter

	// Middleware returns an HTTP middleware configured for this capturer.
	Middleware() Middleware
}
