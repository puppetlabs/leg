// Package norm provides normalization routines for making arbitrary text
// compatible with Kubernetes's field requirements.
package norm

import "strings"

func mapper(ext string) func(rune) rune {
	return func(r rune) rune {
		switch {
		case r >= '0' && r <= '9':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r | 0x20
		case strings.ContainsRune(ext, r):
			return r
		default:
			return '-'
		}
	}
}

var (
	dnsLabelNameMapper  = mapper("")
	qualifiedNameMapper = mapper("._")
)

type AnyDNSLabelNameOptions struct {
	DNSNameRFC DNSNameRFC
}

type AnyDNSLabelNameOption interface {
	ApplyToAnyDNSLabelNameOptions(target *AnyDNSLabelNameOptions)
}

func (o *AnyDNSLabelNameOptions) ApplyOptions(opts []AnyDNSLabelNameOption) {
	for _, opt := range opts {
		opt.ApplyToAnyDNSLabelNameOptions(o)
	}
}

// AnyDNSLabelName normalizes a raw string so that it conforms to the structure
// of a DNS label: it must be all lowercase, contain only alphanumeric
// characters and dashes; and the last character must be alphanumeric. In the
// case of RFC 1035, the first character must be alphabetical, while in the case
// of RFC 1123, the first character must be alphanumeric.
func AnyDNSLabelName(raw string, opts ...AnyDNSLabelNameOption) string {
	o := &AnyDNSLabelNameOptions{
		DNSNameRFC: DNSNameRFC1035,
	}
	o.ApplyOptions(opts)

	mapped := strings.Map(dnsLabelNameMapper, raw)
	if o.DNSNameRFC == DNSNameRFC1123 {
		mapped = strings.TrimLeft(mapped, "-")
	} else {
		mapped = strings.TrimLeft(mapped, "0123456789-")
	}
	if len(mapped) > 63 {
		mapped = mapped[:63]
	}
	mapped = strings.TrimRight(mapped, "-")
	return mapped
}

// AnyDNSLabelNameSuffixed calls AnyDNSLabelName ensuring that the entirety of
// the given suffix is retained if possible.
func AnyDNSLabelNameSuffixed(prefix, suffix string) string {
	remaining := 63 - len(suffix)
	if remaining <= 0 {
		return AnyDNSLabelName(suffix)
	}
	if remaining < len(prefix) {
		prefix = prefix[:remaining]
	}
	return AnyDNSLabelName(prefix + suffix)
}

// AnyDNSSubdomainName normalizes a raw string so that it conforms to the
// structure of a DNS domain (RFC 1123): it must be all lowercase; contain only
// alphanumeric characters and dashes; and the first and last characters must be
// alphanumeric.
func AnyDNSSubdomainName(raw string) string {
	parts := strings.Split(raw, ".")

	var keep int
	for _, part := range parts {
		if part == "" {
			continue
		}

		parts[keep] = AnyDNSLabelName(part, WithDNSNameRFC(DNSNameRFC1123))
		keep++
	}

	mapped := strings.Join(parts[:keep], ".")
	if len(mapped) > 253 {
		mapped = mapped[:253]
	}
	mapped = strings.TrimRight(mapped, ".-")
	return mapped
}

// AnyQualifiedName normalizes a raw string so that it conforms to the structure
// of a Kubernetes qualified name.
//
// Qualified names optionally start with a DNS subdomain followed by a slash.
// They always begin and end with an alphanumeric character, and internally may
// contain alphanumeric characters, dashes, underscores, and dots up to 63
// characters.
func AnyQualifiedName(raw string) string {
	parts := strings.SplitN(raw, "/", 2)
	prefix, name := "", parts[0]
	switch len(parts) {
	case 2:
		prefix, name = AnyDNSSubdomainName(parts[0])+"/", parts[1]
		fallthrough
	case 1:
		name = strings.Map(qualifiedNameMapper, name)
		name = strings.TrimLeft(name, "._-")
		if len(name) > 63 {
			name = name[:63]
		}
		name = strings.TrimRight(name, "._-")
	}
	return prefix + name
}

// MetaName normalizes a Kubernetes metadata name field.
func MetaName(raw string) string {
	return AnyDNSLabelName(raw)
}

// MetaNameSuffixed normalizes a Kubernetes metadata name field ensuring that
// the entirety of the given suffix is retained if possible.
func MetaNameSuffixed(prefix, suffix string) string {
	return AnyDNSLabelNameSuffixed(prefix, suffix)
}

// MetaGenerateName normalizes a Kubernetes metadata generateName field. It is
// opinionated in that it also forces a dash before the five characters
// generated by the Kubernetes API.
func MetaGenerateName(raw string) string {
	mapped := strings.Map(dnsLabelNameMapper, raw)
	mapped = strings.Trim(mapped, "-")
	if len(mapped) > 57 {
		mapped = mapped[:57] + "-"
	} else if mapped[len(mapped)-1] != '-' {
		mapped += "-"
	}
	return mapped
}
