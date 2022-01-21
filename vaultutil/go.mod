module github.com/puppetlabs/leg/vaultutil

go 1.16

require (
	github.com/hashicorp/vault/api v1.2.0
	github.com/hashicorp/vault/sdk v0.2.2-0.20211209191057-0400d2f3deaf // indirect
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	k8s.io/api => k8s.io/api v0.19.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.7
	k8s.io/client-go => k8s.io/client-go v0.19.7
)
