package iso8601

import (
	"fmt"
	"strings"
	"time"

	"github.com/reflect/xparse/xtime"
)

var (
	timeTree = xtime.Compile([]string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z0700",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999Z07:00",
		"2006-01-02",
		"15:04:05",
		"15:04:05Z0700",
		"15:04:05Z07:00",
		"15:04:05.999999999Z07:00",
		"15:04:05.999Z07:00",
	})
)

// Interval is the type of an ISO 8601 time interval. ISO 8601 defines a time
// interval as part of the time axis limited by two instants.
//
// Intervals are represented in one of three ways: by a start and end time, by a
// start time and a duration, or by a duration and an end time. A fourth
// representation specified in ISO 8601, given only by a duration and extended
// context information, is not supported by this implementation.
//
// The two components of an interval are separated by a designator character
// sequence, normally "/". In rare cases, such as when an interval is used in a
// file name, the sequence "--" is used.
//
// For more information, see ISO 8601:2004(E) § 2.1.3, 3.4, and 4.4.
type Interval struct {
	designator string
	start, end time.Time
	duration   Duration
}

// Start returns the absolute start time of this interval regardless of how the
// interval was constructed.
func (i Interval) Start() time.Time {
	if i.start.IsZero() {
		return i.end.Add(i.duration.Before(i.end))
	}

	return i.start
}

// End returns the absolute end time of this interval regardless of how the interval
// was constructed.
func (i Interval) End() time.Time {
	if i.end.IsZero() {
		return i.start.Add(i.duration.After(i.start))
	}

	return i.end
}

// Duration returns the absolute duration of this interval regardless of how the
// interval was constructed.
func (i Interval) Duration() time.Duration {
	return i.End().Sub(i.Start())
}

// Offset returns a new interval offset a given number of times from the current
// interval. It correctly accounts for duration precision. For example,
// offsetting an interval specified as "2018-10-01/P1M" by 2 will return the
// interval specified by "2018-12-01/P1M".
func (i Interval) Offset(n int) (in Interval) {
	in.designator = i.designator

	if i.duration.IsZero() {
		step := time.Duration(n) * i.Duration()

		in.start = i.start.Add(step)
		in.end = i.end.Add(step)
	} else {
		step := i.duration.Mul(n).After(i.Start())

		in.start = i.Start().Add(step)
		in.duration = i.duration
	}

	return
}

// Designator returns the designator character sequence used to separate the
// components of this interval.
func (i Interval) Designator() string {
	if i.designator == "" {
		return "/"
	}

	return i.designator
}

func (i Interval) String() string {
	if !i.duration.IsZero() {
		if !i.start.IsZero() {
			return fmt.Sprintf("%s%s%s", i.start.Format(time.RFC3339), i.Designator(), i.duration)
		} else if !i.end.IsZero() {
			return fmt.Sprintf("%s%s%s", i.duration, i.Designator(), i.end.Format(time.RFC3339))
		} else {
			return i.duration.String()
		}
	} else {
		return fmt.Sprintf("%s%s%s", i.start.Format(time.RFC3339), i.Designator(), i.end.Format(time.RFC3339))
	}
}

func (i Interval) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

func (i *Interval) parse(s string) error {
	sn := strings.SplitN(s, i.Designator(), 2)
	if len(sn) != 2 {
		return ErrInvalidIntervalFormat
	}

	if strings.HasPrefix(sn[0], "P") {
		var err error

		i.duration, err = ParseDuration(sn[0])
		if err != nil {
			return err
		}

		if strings.HasPrefix(sn[1], "P") {
			return ErrInvalidIntervalFormat
		}

		i.end, err = timeTree.Parse(sn[1])
		if err != nil {
			return err
		}
	} else {
		var err error

		i.start, err = timeTree.Parse(sn[0])
		if err != nil {
			return err
		}

		if strings.HasPrefix(sn[1], "P") {
			i.duration, err = ParseDuration(sn[1])
		} else {
			i.end, err = timeTree.Parse(sn[1])
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interval) UnmarshalText(text []byte) error {
	return i.parse(string(text))
}

// ParseInterval parses the given string according to the ISO 8601 time interval
// representation.
func ParseInterval(text string) (Interval, error) {
	var i Interval
	return i, i.parse(text)
}
