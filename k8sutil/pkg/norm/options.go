package norm

type DNSNameRFC string

const (
	DNSNameRFC1035 DNSNameRFC = "1035"
	DNSNameRFC1123 DNSNameRFC = "1123"
)

type WithDNSNameRFC DNSNameRFC

var _ AnyDNSLabelNameOption = WithDNSNameRFC("")

func (wdnr WithDNSNameRFC) ApplyToAnyDNSLabelNameOptions(target *AnyDNSLabelNameOptions) {
	target.DNSNameRFC = DNSNameRFC(wdnr)
}
