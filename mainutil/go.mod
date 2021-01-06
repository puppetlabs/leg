module github.com/puppetlabs/leg/mainutil

go 1.14

require (
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/leg/lifecycle v0.1.0
	github.com/puppetlabs/leg/logging v0.1.0
)

replace github.com/puppetlabs/leg/lifecycle => ../lifecycle
