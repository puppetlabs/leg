package model

const (
	DefaultVaultKubernetesHost = "https://kubernetes.default.svc"

	VaultRootToken = "root-token"
	VaultUnsealKey = "unseal-key"

	VaultKubernetesHost   = "kubernetes_host"
	VaultKubernetesCACert = "kubernetes_ca_cert"

	VaultPathConfig = "config"
	VaultPathKeys   = "keys"
	VaultPathRole   = "role"

	VaultTokenReviewerJWT = "token_reviewer_jwt"

	VaultSysAuth   = "sys/auth"
	VaultSysMounts = "sys/mounts"
)
const (
	KubernetesSecretDataCACertificate = "ca.crt"
	KubernetesSecretDataToken         = "token"
)

type KubernetesCertificateData struct {
	CACertificate string
	Token         string
}

type VaultAuthAccessor struct {
	Name  string
	Rules string
}

type VaultAuthType string

func (vat VaultAuthType) String() string {
	return string(vat)
}

const (
	VaultAuthTypeJWT        VaultAuthType = "jwt"
	VaultAuthTypeKubernetes VaultAuthType = "kubernetes"
)

type VaultJWTRole struct {
	Name           string
	RoleType       string
	BoundAudiences []string
	UserClaim      string
	TokenType      string
	TokenPolicies  []string
	ClaimMappings  map[string]string
}

type VaultKeys struct {
	RootToken  string
	UnsealKeys []string
}

type VaultKubernetesRole struct {
	BoundServiceAccountNames      []string
	BoundServiceAccountNamespaces []string
	Name                          string
	Policies                      []string
	TTL                           string
}

type VaultPolicy struct {
	Name  string
	Rules string
}

type VaultSecretEngine struct {
	Name string
	Type VaultSecretEngineType
}

type VaultSecretEngineType string

func (vset VaultSecretEngineType) String() string {
	return string(vset)
}

const (
	VaultSecretEngineTypeKVV2    VaultSecretEngineType = "kv-v2"
	VaultSecretEngineTypeTransit VaultSecretEngineType = "transit"
)
