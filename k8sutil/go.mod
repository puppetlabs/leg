module github.com/puppetlabs/leg/k8sutil

go 1.14

require (
	github.com/google/uuid v1.1.2
	github.com/puppetlabs/leg/lifecycle v0.0.0-00010101000000-000000000000
	github.com/puppetlabs/leg/timeutil v0.1.0
	github.com/rancher/remotedialer v0.2.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
	sigs.k8s.io/controller-runtime v0.7.0
)

replace github.com/puppetlabs/leg/timeutil => ../timeutil

replace github.com/puppetlabs/leg/mathutil => ../mathutil

replace github.com/puppetlabs/leg/lifecycle => ../lifecycle
