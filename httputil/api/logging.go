package api

import (
	"context"

	logging "github.com/puppetlabs/insights-logging"
)

var (
	logger = logging.Builder().At("insights-stdlib", "httputil", "api")
)

func log(ctx context.Context) logging.Logger {
	return logger.With(ctx).Build()
}
