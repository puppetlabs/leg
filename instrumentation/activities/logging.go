package activities

import (
	"context"

	"github.com/puppetlabs/horsehead/v2/logging"
)

var (
	defaultLogger = logging.Builder().At("horsehead", "instrumentation", "activities")
)

func log(ctx context.Context) logging.Logger {
	return defaultLogger.With(ctx).Build()
}
