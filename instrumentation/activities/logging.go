package activities

import (
	"context"

	"github.com/puppetlabs/leg/logging"
)

var (
	defaultLogger = logging.Builder().At("horsehead", "instrumentation", "activities")
)

func log(ctx context.Context) logging.Logger {
	return defaultLogger.With(ctx).Build()
}
