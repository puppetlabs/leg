package iso8601

import "errors"

var (
	ErrInvalidDurationFormat          = errors.New("iso8601: invalid duration format")
	ErrInvalidIntervalFormat          = errors.New("iso8601: invalid interval format")
	ErrInvalidRecurringIntervalFormat = errors.New("iso8601: invalid recurring interval format")
)
