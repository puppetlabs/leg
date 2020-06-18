package activities

import "errors"

var (
	ErrReportTimeout    = errors.New("activities: report timeout")
	ErrReporterShutdown = errors.New("activities: reporter is shutting down")
)
