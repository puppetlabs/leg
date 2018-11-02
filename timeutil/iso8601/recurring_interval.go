package iso8601

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RecurringInterval is the type of an ISO 8601 recurring time interval. ISO
// 8601 defines a recurring time interval as a series of consecutive time
// intervals of the same duration or nominal duration.
//
// Recurring intervals are represented either by "R<n>/<interval>", where <n>
// indicates the number of repetitions, or by "R/<interval>" to indicate
// unbounded repetition. The first repetition of the recurring interval occurs
// at the start time of the given interval, and it repeats immediately after the
// duration of the interval has elapsed for the number of given repetitions.
//
// The two components of a recurring interval are separated by a designator
// character sequence, normally "/". In rare cases, such as when a recurring
// interval is used in a file name, the sequence "--" is used.
//
// For more information, see ISO 8601:2004(E) § 2.1.17, 3.4, and 4.5.
type RecurringInterval struct {
	designator  string
	repetitions int
	interval    Interval
}

// Repetitions returns the number of repetitions specified by this recurring
// interval. If this recurring interval is unbounded, it sets its second return
// value to false.
func (ri RecurringInterval) Repetitions() (n int, ok bool) {
	return ri.repetitions, ri.repetitions > 0
}

// Interval returns the initial interval of this recurring interval.
func (ri RecurringInterval) Interval() Interval {
	return ri.interval
}

// Next returns the interval starting at or after the given time. If there are
// no intervals after the given time, it sets its second return value to false.
func (ri RecurringInterval) Next(from time.Time) (Interval, bool) {
	// Early check: if the requested time is before the first interval, just
	// jump out.
	start := ri.interval.Start()
	if start.After(from) {
		return ri.interval, true
	}

	// This algorithm is surprising because the actual length of a duration
	// isn't really known ahead of time. It only makes sense relative to a start
	// or end time. To know the "next" value after a whole bunch of iterations,
	// you need to actually converge on what the previous iterations are first.
	//
	// When we let cur := start.Offset(n) such that cur.Start() < from and
	// cur.End() > from, then we know that start.Offset(n+1) will be the next
	// interval.

	duration := ri.interval.Duration()
	if duration <= 0 {
		// This shouldn't ever really happen.
		return ri.interval, true
	}

	n := int(from.Sub(start) / duration)
	for {
		cur := ri.interval.Offset(n)
		start, end := cur.Start(), cur.End()

		if start.After(from) {
			rel := int(start.Sub(from) / cur.Duration())
			if rel < 1 {
				rel = 1
			}

			n -= rel
		} else if end.Before(from) {
			rel := int(from.Sub(end) / cur.Duration())
			if rel < 1 {
				rel = 1
			}

			n += rel
		} else {
			if start != from {
				n++
			}

			break
		}
	}

	// NB: Repetitions is 1-indexed; this function is 0-indexed.
	if ri.repetitions > 0 && n >= ri.repetitions {
		return Interval{}, false
	}

	return ri.interval.Offset(n), true
}

// Designator returns the designator character sequence used to separate the
// components of this recurring interval.
func (ri RecurringInterval) Designator() string {
	if ri.designator == "" {
		return "/"
	}

	return ri.designator
}

func (ri RecurringInterval) String() string {
	it := ri.interval.String()

	if ri.repetitions <= 0 {
		return fmt.Sprintf("R%s%s", ri.Designator(), it)
	}

	return fmt.Sprintf("R%d%s%s", ri.repetitions, ri.Designator(), it)
}

func (ri RecurringInterval) MarshalText() ([]byte, error) {
	return []byte(ri.String()), nil
}

func (ri *RecurringInterval) parseInitial(s string) (rest string, err error) {
	if len(s) == 0 {
		err = ErrInvalidRecurringIntervalFormat
		return
	} else if s[0] != 'R' {
		err = ErrInvalidRecurringIntervalFormat
		return
	}

	return s[1:], nil
}

func (ri *RecurringInterval) parseRepetitions(s string) (rest string, err error) {
	var buf bytes.Buffer

	for i := 0; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		buf.WriteByte(s[i])
	}

	if len := buf.Len(); len > 0 {
		repetitions, err := strconv.ParseInt(buf.String(), 10, 32)
		if err != nil {
			return s, err
		}

		ri.repetitions = int(repetitions)
		return s[len:], nil
	}

	return s, nil
}

func (ri *RecurringInterval) parseDesignator(s string) (rest string, err error) {
	designator := ri.Designator()

	if !strings.HasPrefix(s, designator) {
		err = ErrInvalidRecurringIntervalFormat
		return
	}

	return s[len(designator):], nil
}

func (ri *RecurringInterval) parseInterval(s string) (rest string, err error) {
	ri.interval.designator = ri.designator
	err = ri.interval.parse(s)

	return
}

func (ri *RecurringInterval) parse(s string) error {
	rest, err := parse(s, []parseFunc{
		ri.parseInitial,
		ri.parseRepetitions,
		ri.parseDesignator,
		ri.parseInterval,
	})
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return ErrInvalidRecurringIntervalFormat
	}

	return nil
}

func (ri *RecurringInterval) UnmarshalText(text []byte) error {
	return ri.parse(string(text))
}

// ParseRecurringInterval parses the given string according to the ISO 8601
// recurring time interval representation.
func ParseRecurringInterval(text string) (RecurringInterval, error) {
	var ri RecurringInterval
	return ri, ri.parse(text)
}
