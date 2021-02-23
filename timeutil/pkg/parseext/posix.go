package parseext

import (
	"bytes"
	"time"
)

const (
	MaximumPOSIXBranches = 512

	posixStateStr = iota
	posixStateSpecifier
)

type posixBuffers []*bytes.Buffer

func (p *posixBuffers) WriteString(ss ...string) {
	switch len(ss) {
	case 0:
	case 1:
		for _, buf := range *p {
			buf.WriteString(ss[0])
		}
	default:
		pn := make(posixBuffers, len(*p)*len(ss))
		for i, buf := range *p {
			b := buf.Bytes()

			for j, s := range ss {
				idx := i*len(ss) + j

				pn[idx] = bytes.NewBuffer(make([]byte, 0, len(b)))
				pn[idx].Write(b)
				pn[idx].WriteString(s)
			}
		}

		*p = pn
	}
}

func (p posixBuffers) Strings() []string {
	r := make([]string, len(p))
	for i, buf := range p {
		r[i] = buf.String()
	}

	return r
}

//nolint:gocyclo // It's theoretically possible to express this as a map, but
//               // actually harder to read because of the duplication of keys.
func posixTimeFormat(c rune) ([]string, error) {
	switch c {
	case 'a', 'A':
		return []string{"Mon", "Monday"}, nil
	case 'b', 'B', 'h':
		return []string{"Jan", "January"}, nil
	case 'C':
		// The century of the date (e.g. 19 or 20). Very strange.
		return nil, ErrNotImplemented
	case 'd', 'e':
		return []string{"2", "02"}, nil
	case 'f':
		return []string{"999999999"}, nil
	case 'D':
		return []string{"1/2/06"}, nil
	case 'H':
		return []string{"15"}, nil
	case 'I':
		return []string{"3", "03"}, nil
	case 'j':
		// Day number of year.
		return nil, ErrNotImplemented
	case 'm':
		return []string{"1", "01"}, nil
	case 'M':
		return []string{"4", "04"}, nil
	case 'n', 't':
		return []string{" ", "\n", "\r", "\t", "\v"}, nil
	case 'p':
		return []string{"PM"}, nil
	case 'r':
		return []string{"3:4:5PM"}, nil
	case 'R':
		return []string{"15:4"}, nil
	case 'S':
		return []string{"5", "05"}, nil
	case 'T':
		return []string{"15:4:5"}, nil
	case 'U':
		// Week number of year (first day of the week is Sunday).
		return nil, ErrNotImplemented
	case 'w':
		// Day number of week.
		return nil, ErrNotImplemented
	case 'W':
		// Week number of year (first day of the week is Monday).
		return nil, ErrNotImplemented
	case 'c', 'x', 'X', 'E', 'O':
		// Locale-specific formats.
		return nil, ErrNotImplemented
	case 'y':
		return []string{"06"}, nil
	case 'Y':
		return []string{"2006"}, nil
	case 'z':
		// Time zone. Not supported by POSIX specification, but
		// supported by Python. Seems fairly common.
		return []string{"Z0700", "Z07:00"}, nil
	case '%':
		return []string{"%"}, nil
	default:
		return nil, ErrInvalidFormatString
	}
}

// FromPOSIX converts a POSIX time format string like %Y-%m-%d and converts it
// to a list of time formats compatible with this package.
func FromPOSIX(format string) ([]string, error) {
	bufs := &posixBuffers{&bytes.Buffer{}}
	state := posixStateStr

	for _, c := range format {
		switch state {
		case posixStateStr:
			switch c {
			case '%':
				state = posixStateSpecifier
			default:
				bufs.WriteString(string(c))
			}
		case posixStateSpecifier:
			fmt, err := posixTimeFormat(c)
			if err != nil {
				return nil, err
			}

			bufs.WriteString(fmt...)
			state = posixStateStr
		}

		if len(*bufs) > MaximumPOSIXBranches {
			return nil, ErrFormatStringTooComplex
		}
	}

	return bufs.Strings(), nil
}

// ParsePOSIX parses the given time string according to the given POSIX format.
func ParsePOSIX(format, value string) (time.Time, error) {
	fs, err := FromPOSIX(format)
	if err != nil {
		return time.Time{}, err
	}

	return Compile(fs).Parse(value)
}
