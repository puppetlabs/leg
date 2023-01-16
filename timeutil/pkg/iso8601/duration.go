package iso8601

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"time"
)

// Duration is the type of an ISO 8601 duration. A duration is defined as a
// non-negative quantity attributed to a time interval, the value of which is
// equal to the difference between the time points of the final instant and the
// initial instant of the time interval, when the time points are quantitative
// marks.
//
// "P" is used as a duration indicator. An ISO 8601 representation of a duration
// always starts with a "P" character designator. It is followed by a string in
// the form "<n>W" or "<n>Y<n>M<n>DT<n>H<n>M<n>S".
//
// When a component designator is zero, it may be omitted, except that at least
// one component must be present. The smallest component provided may be a
// fractional span represented as a decimal (for example, PT1.5M).
//
// A third form, P<date>T<time>, is defined, but is not supported by this
// implementation.
//
// For more information, see ISO 8601:2004(E) § 2.1.6, 2.1.7, 3.4, and 4.4.3.
type Duration struct {
	weeks                   *big.Rat
	years, months, days     *big.Rat
	hours, minutes, seconds *big.Rat
}

func (d Duration) get(from time.Time, sign int) time.Duration {
	to := from

	if d.weeks != nil {
		days := big.NewRat(7, 1)
		days.Mul(days, d.weeks)

		to, _ = addRatDays(to, days)
		return to.Sub(from)
	}

	steps := []struct {
		x  *big.Rat
		fn func(t time.Time, x *big.Rat) (time.Time, bool)
	}{
		{d.years, addRatYears},
		{d.months, addRatMonths},
		{d.days, addRatDays},
		{d.hours, addRatHours},
		{d.minutes, addRatMinutes},
		{d.seconds, addRatSeconds},
	}
	for _, step := range steps {
		if step.x == nil {
			continue
		}

		x := step.x
		if sign < 0 {
			x = (&big.Rat{}).Neg(step.x)
		}

		var cont bool
		to, cont = step.fn(to, x)
		if !cont {
			break
		}
	}

	return to.Sub(from)
}

// Before returns an absolute duration of time in Go's native duration format
// (nanoseconds) as if the duration occurred immediately before the given time.
func (d Duration) Before(from time.Time) time.Duration {
	return d.get(from, -1)
}

// After returns an absolute duration of time in Go's native duration format
// (nanoseconds) as if the duration occurred immediately after the given time.
func (d Duration) After(from time.Time) time.Duration {
	return d.get(from, 1)
}

// Mul returns a new duration given by the product of the current duration and
// the given operand. It is precise to the definition of the current duration;
// for example, multiplying a span of 1 year 6 months by 2 will return 3 years
// exactly.
func (d Duration) Mul(n int) (dn Duration) {
	muls := []struct {
		in  *big.Rat
		out **big.Rat
	}{
		{d.years, &dn.years},
		{d.months, &dn.months},
		{d.days, &dn.days},
		{d.hours, &dn.hours},
		{d.minutes, &dn.minutes},
		{d.seconds, &dn.seconds},
	}
	for _, mul := range muls {
		if mul.in == nil {
			continue
		}

		*mul.out = (&big.Rat{}).Mul(mul.in, big.NewRat(int64(n), 1))
	}

	return
}

// IsZero returns true if this duration has not been set. It does not return
// zero if an explicit duration of no time has been provided.
func (d Duration) IsZero() bool {
	return d.weeks == nil &&
		d.years == nil && d.months == nil && d.days == nil &&
		d.hours == nil && d.minutes == nil && d.seconds == nil
}

func (d Duration) String() string {
	buf := bytes.NewBufferString("P")

	if d.weeks != nil {
		buf.WriteString(d.weeks.FloatString(-1))
		buf.WriteByte('W')
	} else {
		if d.years != nil {
			buf.WriteString(d.years.FloatString(-1))
			buf.WriteByte('Y')
		}

		if d.months != nil {
			buf.WriteString(d.months.FloatString(-1))
			buf.WriteByte('M')
		}

		if d.days != nil {
			buf.WriteString(d.days.FloatString(-1))
			buf.WriteByte('D')
		}

		if d.hours != nil || d.minutes != nil || d.seconds != nil {
			buf.WriteByte('T')
		}

		if d.hours != nil {
			buf.WriteString(d.hours.FloatString(-1))
			buf.WriteByte('H')
		}

		if d.minutes != nil {
			buf.WriteString(d.minutes.FloatString(-1))
			buf.WriteByte('M')
		}

		if d.seconds != nil {
			buf.WriteString(d.seconds.FloatString(-1))
			buf.WriteByte('S')
		}
	}

	return buf.String()
}

func (d *Duration) parseInitial(s string) (rest string, err error) {
	if s == "" {
		err = io.EOF
		return
	} else if s[0] != 'P' {
		err = ErrInvalidDurationFormat
		return
	}

	return s[1:], nil
}

func (d *Duration) parseElement(id rune, target **big.Rat) parseFunc {
	return func(s string) (string, error) {
		buf := bytes.NewBufferString(s)

		var candidate rune
		var value big.Rat

		n, err := fmt.Fscanf(buf, "%f%c", &value, &candidate)
		if err != nil || n != 2 || candidate != id {
			return s, nil
		}

		*target = &value
		return buf.String(), nil
	}
}

func (d *Duration) parseWeek(s string) (rest string, err error) {
	return d.parseElement('W', &d.weeks)(s)
}

func (d *Duration) parseDate(s string) (rest string, err error) {
	if s == "" {
		return
	} else if d.weeks != nil {
		err = ErrInvalidDurationFormat
		return
	}

	return parse(s, []parseFunc{
		d.parseElement('Y', &d.years),
		d.parseElement('M', &d.months),
		d.parseElement('D', &d.days),
	})
}

func (d *Duration) parseTime(s string) (rest string, err error) {
	if s == "" {
		if d.IsZero() {
			// We must have a date or a time. Otherwise not to spec.
			err = ErrInvalidDurationFormat
			return
		}

		return
	} else if s[0] != 'T' {
		return
	} else if d.weeks != nil {
		// We can't have a week spec with a time spec.
		err = ErrInvalidDurationFormat
		return
	}

	return parse(s[1:], []parseFunc{
		d.parseElement('H', &d.hours),
		d.parseElement('M', &d.minutes),
		d.parseElement('S', &d.seconds),
	})
}

func (d *Duration) parse(s string) error {
	rest, err := parse(s, []parseFunc{
		d.parseInitial,
		d.parseWeek,
		d.parseDate,
		d.parseTime,
	})
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return ErrInvalidDurationFormat
	}

	return nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	return d.parse(string(text))
}

// ParseDuration parses the given string according to the ISO 8601 duration
// representation.
func ParseDuration(text string) (d Duration, err error) {
	err = d.parse(text)
	return
}
