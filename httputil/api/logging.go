package api

import (
	"context"

	logging "github.com/puppetlabs/leg/logging"
)

var (
	logger = logging.Builder().At("leg", "httputil", "api")
)

func log(ctx context.Context) logging.Logger {
	return logger.With(ctx).Build()
}
