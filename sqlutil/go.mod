module github.com/puppetlabs/leg/sqlutil

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/puppetlabs/leg/lifecycle v0.1.0
	github.com/stretchr/testify v1.6.1
)

replace github.com/puppetlabs/leg/lifecycle => ../lifecycle
