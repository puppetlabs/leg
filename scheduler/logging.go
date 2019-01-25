package scheduler

import (
	"context"

	logging "github.com/puppetlabs/insights-logging"
)

var (
	logger = logging.Builder().At("insights-stdlib", "scheduler")
)

func log(ctx context.Context) logging.Logger {
	return logger.With(ctx).Build()
}
