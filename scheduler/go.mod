module github.com/puppetlabs/leg/scheduler

go 1.14

require (
	github.com/puppetlabs/errawr-gen v1.0.1
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/leg/instrumentation v0.1.3
	github.com/puppetlabs/leg/logging v0.1.0
	github.com/puppetlabs/leg/netutil v0.1.0
	github.com/puppetlabs/leg/request v0.1.0
	github.com/stretchr/testify v1.6.1
)

replace github.com/puppetlabs/leg/instrumentation => ../instrumentation
