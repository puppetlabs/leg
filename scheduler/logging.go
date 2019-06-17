package scheduler

import (
	"context"

	logging "github.com/puppetlabs/insights-logging"
)

var (
	logger = logging.Builder().At("horsehead", "scheduler")
)

func log(ctx context.Context) logging.Logger {
	return logger.With(ctx).Build()
}
