package alertstest

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/puppetlabs/horsehead/v2/instrumentation/alerts/trackers"
)

type Capturer struct {
	tags []trackers.Tag

	ReporterRecorders chan *ReporterRecorder
}

func (c Capturer) WithNewTrace() trackers.Capturer {
	return c
}

func (c Capturer) WithAppPackages(packages []string) trackers.Capturer {
	return c
}

func (c Capturer) WithUser(u trackers.User) trackers.Capturer {
	return c
}

func (c Capturer) WithTags(tags ...trackers.Tag) trackers.Capturer {
	return &Capturer{
		tags:              append(append([]trackers.Tag{}, c.tags...), tags...),
		ReporterRecorders: c.ReporterRecorders,
	}
}

func (c Capturer) Try(ctx context.Context, fn func(ctx context.Context)) (rv interface{}) {
	defer func() {
		rv = recover()
		if nil != rv {
			debug.PrintStack()
		}
	}()

	fn(ctx)
	return nil
}

func (c Capturer) Capture(err error) trackers.Reporter {
	rr := &ReporterRecorder{
		err:  err,
		tags: c.tags,
	}

	c.ReporterRecorders <- rr

	return rr
}

func (c Capturer) CaptureMessage(message string) trackers.Reporter {
	return c.Capture(fmt.Errorf(message))
}

func (c Capturer) Middleware() trackers.Middleware {
	return &Middleware{c: &c}
}

func NewCapturer() *Capturer {
	return &Capturer{
		ReporterRecorders: make(chan *ReporterRecorder),
	}
}
