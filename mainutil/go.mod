module github.com/puppetlabs/leg/mainutil

go 1.14

require (
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/leg/lifecycle v0.0.0-00010101000000-000000000000
	github.com/puppetlabs/leg/logging v0.0.0-00010101000000-000000000000
)

replace github.com/puppetlabs/leg/lifecycle => ../lifecycle

replace github.com/puppetlabs/leg/logging => ../logging
