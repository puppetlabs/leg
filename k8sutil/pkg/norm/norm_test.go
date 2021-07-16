package norm_test

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/puppetlabs/leg/k8sutil/pkg/norm"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation"
)

func TestNorm(t *testing.T) {
	tests := []struct {
		Func      func(string) string
		Validator func(string) []string
		Raw       string
		Expected  string
	}{
		{
			Func: func(raw string) string {
				return norm.AnyDNSLabelName(raw)
			},
			Validator: validation.IsDNS1035Label,
			Raw:       "foo",
			Expected:  "foo",
		},
		{
			Func: func(raw string) string {
				return norm.AnyDNSLabelName(raw)
			},
			Validator: validation.IsDNS1035Label,
			Raw:       strings.Repeat("foo", 100),
			Expected:  strings.Repeat("foo", 100)[:63],
		},
		{
			Func:      norm.AnyDNSSubdomainName,
			Validator: validation.IsDNS1123Subdomain,
			Raw:       strings.Repeat("foo.", 100),
			Expected:  strings.Repeat("foo.", 100)[:253],
		},
		{
			Func: func(raw string) string {
				return norm.AnyDNSLabelName(raw)
			},
			Validator: validation.IsDNS1035Label,
			Raw:       "$sFj.Mj-%29A&zKL",
			Expected:  "sfj-mj--29a-zkl",
		},
		{
			Func: func(raw string) string {
				return norm.AnyDNSLabelName(raw)
			},
			Validator: validation.IsDNS1035Label,
			Raw:       "0123456789abcd012345",
			Expected:  "abcd012345",
		},
		{
			Func: func(raw string) string {
				return norm.AnyDNSLabelName(raw, norm.WithDNSNameRFC(norm.DNSNameRFC1123))
			},
			Validator: validation.IsDNS1123Label,
			Raw:       "0123456789abcd012345",
			Expected:  "0123456789abcd012345",
		},
		{
			Func: func(raw string) string {
				return norm.AnyDNSLabelName(raw)
			},
			Validator: validation.IsDNS1035Label,
			Raw:       strings.Repeat("x", 62) + "-a",
			Expected:  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
		{
			Func:      norm.AnyQualifiedName,
			Validator: validation.IsQualifiedName,
			Raw:       "foo",
			Expected:  "foo",
		},
		{
			Func:      norm.AnyQualifiedName,
			Validator: validation.IsQualifiedName,
			Raw:       "foo.example.com/bar",
			Expected:  "foo.example.com/bar",
		},
		{
			Func:      norm.AnyQualifiedName,
			Validator: validation.IsQualifiedName,
			Raw:       "foo.example.com/$sFj.Mj-_29A&zKL",
			Expected:  "foo.example.com/sfj.mj-_29a-zkl",
		},
		{
			Func:      norm.AnyQualifiedName,
			Validator: validation.IsQualifiedName,
			Raw:       "foo.bar$..example.com^/bar",
			Expected:  "foo.bar.example.com/bar",
		},
		{
			Func:      norm.AnyDNSSubdomainName,
			Validator: validation.IsDNS1123Subdomain,
			Raw:       "$sFj.Mj-%29A&zKL",
			Expected:  "sfj.mj--29a-zkl",
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
			normalized := test.Func(test.Raw)
			if test.Validator != nil {
				assert.Empty(t, test.Validator(normalized))
			}
			assert.Equal(t, test.Expected, normalized)
		})
	}
}

func TestNormSuffixed(t *testing.T) {
	tests := []struct {
		Func           func(string, string) string
		Validator      func(string) []string
		Prefix, Suffix string
		Expected       string
	}{
		{
			Func:      norm.AnyDNSLabelNameSuffixed,
			Validator: validation.IsDNS1035Label,
			Prefix:    "foo",
			Suffix:    "bar",
			Expected:  "foobar",
		},
		{
			Func:      norm.AnyDNSLabelNameSuffixed,
			Validator: validation.IsDNS1035Label,
			Prefix:    strings.Repeat("foo", 100),
			Suffix:    "bar",
			Expected:  strings.Repeat("foo", 20) + "bar",
		},
		{
			Func:      norm.AnyDNSLabelNameSuffixed,
			Validator: validation.IsDNS1035Label,
			Prefix:    strings.Repeat("foo", 100),
			Suffix:    strings.Repeat("bar", 100),
			Expected:  strings.Repeat("bar", 21),
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%s(%q, %q)", runtime.FuncForPC(reflect.ValueOf(test.Func).Pointer()).Name(), test.Prefix, test.Suffix)
		t.Run(name, func(t *testing.T) {
			normalized := test.Func(test.Prefix, test.Suffix)
			if test.Validator != nil {
				assert.Empty(t, test.Validator(normalized))
			}
			assert.Equal(t, test.Expected, normalized)
		})
	}
}
