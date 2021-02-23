package parseext_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/parseext"
	"github.com/stretchr/testify/assert"
)

func TestFromPOSIX(t *testing.T) {
	valid := []struct {
		Format    string
		Candidate string
		Expected  string
	}{
		{`%Y-%m-%dT%H:%M:%S%z`, "2017-04-13T00:31:00-0400", "2017-04-13T00:31:00-04:00"},
		{`%Y-%m-%dT%H:%M:%S%z`, "2017-04-13T00:31:00Z", "2017-04-13T00:31:00Z"},
		{`%Y-%m-%dT%H:%M:%S.%f%z`, "2017-04-13T00:31:00.123456Z", "2017-04-13T00:31:00Z"},
		{`%a, %b %d, %Y %r %z`, "Thu, Apr 13, 2017 00:31:00AM -0400", "2017-04-13T00:31:00-04:00"},
		{`%a, %b %d, %Y %r %z`, "Thursday, April 13, 2017 00:31:00AM -04:00", "2017-04-13T00:31:00-04:00"},
		{`%d/%m/%Y`, "05/04/2017", "2017-04-05T00:00:00Z"},
		{`%d/%m/%Y`, "5/4/2017", "2017-04-05T00:00:00Z"},
		{`%Y%m%d`, "20170405", "2017-04-05T00:00:00Z"},
	}

	for _, v := range valid {
		fs, err := parseext.FromPOSIX(v.Format)
		assert.NoError(t, err)

		p, err := parseext.Compile(fs).Parse(v.Candidate)
		assert.NoError(t, err)

		assert.Equal(t, v.Expected, p.Format(time.RFC3339), "for format %s and candidate %s", v.Format, v.Candidate)
	}
}
