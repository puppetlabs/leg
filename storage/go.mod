module github.com/puppetlabs/leg/storage

go 1.14

replace github.com/puppetlabs/leg/workdir => ../workdir

require (
	cloud.google.com/go/storage v1.12.0
	github.com/google/uuid v1.1.2
	github.com/puppetlabs/leg/workdir v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1
	google.golang.org/api v0.35.0
)
