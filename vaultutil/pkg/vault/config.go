package vault

import "net/url"

type VaultConfig struct {
	JWTAuthPath        string
	JWTMount           string
	JWTPublicKey       string
	KubernetesAuthPath string
	KubernetesMount    string
	Name               string
	Namespace          string
	ServiceAccount     string
	Token              string
	UnsealKey          string
	VaultAddr          *url.URL
}
