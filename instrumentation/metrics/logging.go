package metrics

import (
	"context"

	"github.com/puppetlabs/leg/logging"
)

var (
	defaultLogger = logging.Builder().At("horsehead", "instrumentation", "metrics")
)

func log(ctx context.Context) logging.Logger {
	return defaultLogger.With(ctx).Build()
}
