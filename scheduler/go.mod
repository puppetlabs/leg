module github.com/puppetlabs/leg/scheduler

go 1.14

require (
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/puppetlabs/errawr-gen v1.0.1
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/leg/instrumentation v0.0.0-00010101000000-000000000000
	github.com/puppetlabs/leg/logging v0.0.0-00010101000000-000000000000
	github.com/puppetlabs/leg/netutil v0.0.0-00010101000000-000000000000
	github.com/puppetlabs/leg/request v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1
)

replace github.com/puppetlabs/leg/instrumentation => ../instrumentation

replace github.com/puppetlabs/leg/request => ../request

replace github.com/puppetlabs/leg/netutil => ../netutil

replace github.com/puppetlabs/leg/logging => ../logging

replace github.com/puppetlabs/leg/scheduler => ./
