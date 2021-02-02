package norm_test

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/puppetlabs/leg/k8sutil/pkg/norm"
	"github.com/stretchr/testify/assert"
)

func TestNorm(t *testing.T) {
	tests := []struct {
		Func     func(string) string
		Raw      string
		Expected string
	}{
		{
			Func:     norm.AnyDNSLabelName,
			Raw:      "foo",
			Expected: "foo",
		},
		{
			Func:     norm.AnyDNSLabelName,
			Raw:      strings.Repeat("foo", 100),
			Expected: strings.Repeat("foo", 100)[:63],
		},
		{
			Func:     norm.AnyDNSSubdomainName,
			Raw:      strings.Repeat("foo.", 100),
			Expected: strings.Repeat("foo.", 100)[:253],
		},
		{
			Func:     norm.AnyDNSLabelName,
			Raw:      "$sFj.Mj-%29A&zKL",
			Expected: "sfj-mj--29a-zkl",
		},
		{
			Func:     norm.AnyDNSSubdomainName,
			Raw:      "$sFj.Mj-%29A&zKL",
			Expected: "sfj.mj--29a-zkl",
		},
		{
			Func:     norm.MetaGenerateName,
			Raw:      "foo-",
			Expected: "foo-",
		},
		{
			Func:     norm.MetaGenerateName,
			Raw:      "$sFj.Mj-%29A&zKL",
			Expected: "sfj-mj--29a-zkl-",
		},
		{
			Func:     norm.MetaGenerateName,
			Raw:      strings.Repeat("foo", 100),
			Expected: strings.Repeat("foo", 100)[:57] + "-",
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%s(%q)", runtime.FuncForPC(reflect.ValueOf(test.Func).Pointer()).Name(), test.Raw)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Func(test.Raw))
		})
	}
}

func TestNormSuffixed(t *testing.T) {
	tests := []struct {
		Func           func(string, string) string
		Prefix, Suffix string
		Expected       string
	}{
		{
			Func:     norm.AnyDNSLabelNameSuffixed,
			Prefix:   "foo",
			Suffix:   "bar",
			Expected: "foobar",
		},
		{
			Func:     norm.AnyDNSLabelNameSuffixed,
			Prefix:   strings.Repeat("foo", 100),
			Suffix:   "bar",
			Expected: strings.Repeat("foo", 20) + "bar",
		},
		{
			Func:     norm.AnyDNSLabelNameSuffixed,
			Prefix:   strings.Repeat("foo", 100),
			Suffix:   strings.Repeat("bar", 100),
			Expected: strings.Repeat("bar", 21),
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%s(%q, %q)", runtime.FuncForPC(reflect.ValueOf(test.Func).Pointer()).Name(), test.Prefix, test.Suffix)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Func(test.Prefix, test.Suffix))
		})
	}
}
