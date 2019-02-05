package api_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppetlabs/insights-stdlib/httputil/api"
	"github.com/stretchr/testify/assert"
)

func TestSetContentDispositionHeader(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "test.pdf",
			Expected: "test.pdf",
		},
		{
			Input:    "tést.pdf",
			Expected: "test.pdf",
		},
		{
			Input:    strings.Repeat("N", 1024) + ".txt",
			Expected: strings.Repeat("N", 252) + ".txt",
		},
		{
			Input:    "N." + strings.Repeat("E", 1024),
			Expected: "N." + strings.Repeat("E", 254),
		},
		{
			Input:    strings.Repeat("N", 1024) + "." + strings.Repeat("E", 1024),
			Expected: "N." + strings.Repeat("E", 254),
		},
		{
			Input:    "this\r\nis\r\nẔ̧̻͉̼̺̙͙͖̙̺̼̥͙̗͔̰̲͉ͥ͊͛ͯ̂̽̏̒͆ͤ̑A̴̧̡̪͚̲͉͕͚̳͈̩̯͖͓͑̉́ͤͬ̽L̨̝̩̮̯̤͙͓͔̮̙̺͖̗̝͈̗̰ͬ̍ͮ͌ͨ̋͘Ģ̴͉̺̪̰̼͎̿ͨ̄̅ͯ͆ͥ̇̌̒ͪͫͩ̌̿̃́̏͢͠Ò̷̧ͤ̐͐ͫ͂͐̽ͦͧ̐̓̉ͪ̅̉́͂͛͏̵̤̲̙̝̻̹̼̼͓̰͉̹̫̘̪̹̮̖!̵̯̝̗̟̱̼̪̫̞̲̯̼̝̋ͦ͂̿̌͞ͅ",
			Expected: "this_is_ZALGO_",
		},
		{
			Input:    "\xe2\x82\xa1INVALID.pdf",
			Expected: "_INVALID.pdf",
		},
		{
			Input:    "",
			Expected: "attachment",
		},
		{
			Input:    ".pdf",
			Expected: "attachment.pdf",
		},
	}
	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			res := httptest.NewRecorder()
			api.SetContentDispositionHeader(res, test.Input)

			assert.Equal(t, fmt.Sprintf("attachment; filename=%s", test.Expected), res.Header().Get("content-disposition"))
		})
	}
}
