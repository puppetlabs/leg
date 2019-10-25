package xtime

import (
	"time"
)

func TimezoneOffsetAt(offset string, t time.Time) (hours int, minutes int, err error) {
	loc, err := time.LoadLocation(offset)
	if err != nil {
		return 0, 0, err
	}

	_, diff := t.In(loc).Zone()

	hours = int(diff / 3600)
	minutes = int((diff % 3600) / 60)
	if minutes < 0 {
		minutes = -minutes
	}

	return
}

func TimezoneOffset(offset string) (hours int, minutes int, err error) {
	return TimezoneOffsetAt(offset, time.Now())
}
