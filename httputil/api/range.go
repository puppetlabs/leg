package api

import (
	"fmt"
	"strings"
	"regexp"
	"strconv"
)

type RangeSpec struct {
	First        *int64
	Last         *int64
	// If non-null First/Last are null.
	SuffixLength *int64
}

type RangeHeader struct {
	Unit  string // always "bytes" for now
	Specs []RangeSpec
}

var rangeSpecRegex = regexp.MustCompile(`^\s*(\d*)\s*-\s*(\d*)\s*$`)

func ScanRangeHeader(header string) (*RangeHeader, error) {
	if len(header) == 0 {
		return nil, nil
	}
	eq := strings.Index(header, "=")
	if eq < 1 {
		return nil, fmt.Errorf("Invalid Range header")
	}
	unit := strings.TrimSpace(header[0:eq])
	if "bytes" != unit {
		return nil, fmt.Errorf("Unsupported Range header unit=%s", unit)
	}
	specStrs := strings.Split(header[eq+1:], ",")
	var specs []RangeSpec
	for _, specStr := range specStrs {
		matches := rangeSpecRegex.FindStringSubmatch(specStr)
		if len(matches) <= 0 {
			return nil, fmt.Errorf("Invalid Range header, expected %s to be digits followed by '-' followed by digits", specStr)
		}
		var first, last, suffixLength *int64
		if len(matches[1]) > 0 {
			v, _ := strconv.ParseInt(matches[1], 10, 64)
			first = &v
			if len(matches[2]) > 0 {
				v2, _ := strconv.ParseInt(matches[2], 10, 64)
				if v2 < *first {
					return nil, fmt.Errorf(
						`Unsatisfiable byte range %s`, specStr)
				}
				last = &v2
			}
		} else if len(matches[2]) > 0 {
			v, _ := strconv.ParseInt(matches[2], 10, 64)
			suffixLength = &v
		} else {
			return nil, fmt.Errorf("Invalid Range header, expected more than just a '-'")
		}
		specs = append(specs, RangeSpec {
			First: first,
			Last: last,
			SuffixLength: suffixLength,
		})
			
	}
	return &RangeHeader {
		Unit: unit,
		Specs: specs,
	}, nil
}
