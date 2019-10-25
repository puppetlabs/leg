package xtime

import (
	"bytes"
)

const (
	MaximumPOSIXBranches = 512

	posixStateStr = iota
	posixStateSpecifier
)

type posixBuffers []*bytes.Buffer

func (p *posixBuffers) WriteRune(rs ...rune) {
	switch len(rs) {
	case 0:
	case 1:
		for _, buf := range *p {
			buf.WriteRune(rs[0])
		}
	default:
		pn := make(posixBuffers, len(*p)*len(rs))
		for i, buf := range *p {
			b := buf.Bytes()

			for j, r := range rs {
				idx := i*len(rs) + j

				pn[idx] = bytes.NewBuffer(make([]byte, 0, len(b)))
				pn[idx].Write(b)
				pn[idx].WriteRune(r)
			}
		}

		*p = pn
	}
}

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
		r[i] = string(buf.Bytes())
	}

	return r
}

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
				bufs.WriteRune(c)
			}
		case posixStateSpecifier:
			switch c {
			case 'a', 'A':
				bufs.WriteString("Mon", "Monday")
			case 'b', 'B', 'h':
				bufs.WriteString("Jan", "January")
			case 'C':
				// The century of the date (e.g. 19 or 20). Very strange.
				return nil, ErrNotImplemented
			case 'd', 'e':
				bufs.WriteString("2", "02")
			case 'f':
				bufs.WriteString("999999999")
			case 'D':
				bufs.WriteString("1/2/06")
			case 'H':
				bufs.WriteString("15")
			case 'I':
				bufs.WriteString("3", "03")
			case 'j':
				// Day number of year.
				return nil, ErrNotImplemented
			case 'm':
				bufs.WriteString("1", "01")
			case 'M':
				bufs.WriteString("4", "04")
			case 'n', 't':
				bufs.WriteString(" ", "\n", "\r", "\t", "\v")
			case 'p':
				bufs.WriteString("PM")
			case 'r':
				bufs.WriteString("3:4:5PM")
			case 'R':
				bufs.WriteString("15:4")
			case 'S':
				bufs.WriteString("5", "05")
			case 'T':
				bufs.WriteString("15:4:5")
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
				bufs.WriteString("06")
			case 'Y':
				bufs.WriteString("2006")
			case 'z':
				// Time zone. Not supported by POSIX specification, but
				// supported by Python. Seems fairly common.
				bufs.WriteString("Z0700", "Z07:00")
			case '%':
				bufs.WriteRune('%')
			default:
				return nil, ErrInvalidFormatString
			}

			state = posixStateStr
		}

		if len(*bufs) > MaximumPOSIXBranches {
			return nil, ErrFormatStringTooComplex
		}
	}

	return bufs.Strings(), nil
}
