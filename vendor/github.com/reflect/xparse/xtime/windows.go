//go:generate go run windows_gen.go -output-path windows_cldr_supplement.go

package xtime

import (
	"fmt"
	"time"
)

func WindowsLocationString(l *time.Location) (name string, err error) {
	name, ok := tzdataWindowsMapping[l.String()]
	if !ok {
		err = fmt.Errorf("Unicode CLDR mapping for time zone %q does not exist", l)
	}

	return
}
