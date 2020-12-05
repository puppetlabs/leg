module github.com/puppetlabs/leg/instrumentation

go 1.14

require (
	github.com/aws/aws-sdk-go v1.27.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/getsentry/raven-go v0.2.0
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/prometheus/client_golang v1.8.0
	github.com/puppetlabs/errawr-gen v1.0.1
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/leg/logging v0.1.0
	github.com/puppetlabs/leg/netutil v0.1.0
	github.com/puppetlabs/leg/scheduler v0.1.1
	github.com/segmentio/backo-go v0.0.0-20200129164019-23eae7c10bd3 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	gopkg.in/intercom/intercom-go.v2 v2.0.0-20200217143803-6ffc0627261a
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
)

replace (
	github.com/puppetlabs/leg/logging => ../logging
	github.com/puppetlabs/leg/netutil => ../netutil
	github.com/puppetlabs/leg/scheduler => ../scheduler
)
