module github.com/puppetlabs/leg/sqlutil

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.4.1 // indirect
	github.com/puppetlabs/leg/lifecycle v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2 // indirect
	golang.org/x/sys v0.0.0-20200523222454-059865788121 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20200601175630-2caf76543d99 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/puppetlabs/leg/lifecycle => ../lifecycle
