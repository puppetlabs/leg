package vault

import (
	"context"
	"fmt"
	"path"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/puppetlabs/leg/vaultutil/pkg/model"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VaultSystemManager struct {
	kubeClient  client.Client
	vaultClient *vaultapi.Client
	vaultConfig *VaultConfig
}

func (vsm *VaultSystemManager) Initialize(ctx context.Context) (*model.VaultKeys, error) {
	response, err := vsm.vaultClient.Sys().Init(
		&vaultapi.InitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		},
	)
	if err != nil {
		return nil, err
	}

	vaultKeys := &model.VaultKeys{
		RootToken:  response.RootToken,
		UnsealKeys: response.Keys,
	}

	err = vsm.CreateCredentials(ctx, vaultKeys)
	if err != nil {
		return nil, err
	}

	return vaultKeys, nil
}

func (vsm *VaultSystemManager) GetCredentials(ctx context.Context) (*model.VaultKeys, error) {
	key := types.NamespacedName{
		Name:      vsm.vaultConfig.Name,
		Namespace: vsm.vaultConfig.Namespace,
	}

	secret := &corev1.Secret{}
	err := vsm.kubeClient.Get(ctx, key, secret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &model.VaultKeys{
		RootToken:  string(secret.Data[model.VaultRootToken]),
		UnsealKeys: []string{string(secret.Data[model.VaultUnsealKey])},
	}, nil
}

func (vsm *VaultSystemManager) CreateCredentials(ctx context.Context, vaultKeys *model.VaultKeys) error {
	objectMeta := metav1.ObjectMeta{
		Name:      vsm.vaultConfig.Name,
		Namespace: vsm.vaultConfig.Namespace,
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		StringData: map[string]string{
			model.VaultRootToken: vaultKeys.RootToken,
			model.VaultUnsealKey: vaultKeys.UnsealKeys[0],
		},
	}

	if err := vsm.kubeClient.Create(ctx, secret); err != nil {
		return err
	}

	return nil
}

func (vsm *VaultSystemManager) SetToken(vaultKeys *model.VaultKeys) {
	vsm.vaultClient.SetToken(vaultKeys.RootToken)
}

func (vsm *VaultSystemManager) Unseal(vaultKeys *model.VaultKeys) error {
	_, err := vsm.vaultClient.Sys().UnsealWithOptions(
		&vaultapi.UnsealOpts{
			Key: vaultKeys.UnsealKeys[0],
		})
	if err != nil {
		return err
	}

	return nil
}

func (vsm *VaultSystemManager) RegisterPlugins(plugins []*vaultapi.RegisterPluginInput) error {
	for _, plugin := range plugins {
		err := vsm.vaultClient.Sys().RegisterPlugin(plugin)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vsm *VaultSystemManager) GetAuthMethod(name string) (*vaultapi.AuthMount, error) {
	auth, err := vsm.vaultClient.Sys().ListAuth()
	if err != nil {
		return nil, err
	}

	for n, a := range auth {
		if n == name {
			return a, nil
		}
	}

	return nil, nil
}

func (vsm *VaultSystemManager) PutPolicies(policies []*model.VaultPolicy) error {
	for _, policy := range policies {
		err := vsm.vaultClient.Sys().PutPolicy(policy.Name, policy.Rules)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vsm *VaultSystemManager) ConfigureKubernetesAuthRoles(roles []*model.VaultKubernetesRole) error {
	for _, role := range roles {
		_, err := vsm.vaultClient.Logical().Write(
			path.Join(vsm.vaultConfig.KubernetesAuthPath, model.VaultPathRole, role.Name),
			map[string]interface{}{
				"bound_service_account_names":      role.BoundServiceAccountNames,
				"bound_service_account_namespaces": role.BoundServiceAccountNamespaces,
				"ttl":                              role.TTL,
				"policies":                         role.Policies,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vsm *VaultSystemManager) EnableSecretEngines(secretEngines []*model.VaultSecretEngine) error {
	for _, se := range secretEngines {
		_, err := vsm.vaultClient.Logical().Write(
			path.Join(model.VaultSysMounts, se.Name), vsm.mapSecretEngineType(se))
		if err != nil {
			return err
		}
	}

	return nil
}

func (vsm *VaultSystemManager) CreateTransitKey(mount, key string) error {
	_, err := vsm.vaultClient.Logical().Write(
		path.Join(mount, model.VaultPathKeys, key),
		map[string]interface{}{
			"derived": true,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (vsm *VaultSystemManager) EnableJWTAuth() error {
	authEnabled, err := vsm.IsJWTAuthEnabled()
	if err != nil {
		return err
	}

	if authEnabled {
		return nil
	}

	err = vsm.vaultClient.Sys().EnableAuthWithOptions(vsm.vaultConfig.JWTMount,
		&vaultapi.EnableAuthOptions{
			Type: model.VaultAuthTypeJWT.String(),
		})
	if err != nil {
		return err
	}

	return nil
}

func (vsm *VaultSystemManager) IsJWTAuthEnabled() (bool, error) {
	return vsm.isAuthEnabled(vsm.vaultConfig.JWTMount)
}

func (vsm *VaultSystemManager) ConfigureJWTAuth(ctx context.Context) error {
	if vsm.vaultConfig.JWTPublicKey != "" {
		payload := map[string]interface{}{
			"jwt_supported_algs":     []string{"RS256", "RS512"},
			"jwt_validation_pubkeys": []string{vsm.vaultConfig.JWTPublicKey},
		}

		if _, err := vsm.vaultClient.Logical().Write(
			path.Join(vsm.vaultConfig.JWTAuthPath, model.VaultPathConfig), payload); err != nil {
			return err
		}
	}

	return nil
}

func (vsm *VaultSystemManager) ConfigureJWTAuthRoles(roles []*model.VaultJWTRole) error {
	for _, role := range roles {
		_, err := vsm.vaultClient.Logical().Write(
			path.Join(vsm.vaultConfig.JWTAuthPath, model.VaultPathRole, role.Name),
			map[string]interface{}{
				"bound_audiences": role.BoundAudiences,
				"claim_mappings":  role.ClaimMappings,
				"name":            role.Name,
				"role_type":       role.RoleType,
				"token_policies":  role.TokenPolicies,
				"token_type":      role.TokenType,
				"user_claim":      role.UserClaim,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vsm *VaultSystemManager) EnableKubernetesAuth() error {
	authEnabled, err := vsm.IsKubernetesAuthEnabled()
	if err != nil {
		return err
	}

	if authEnabled {
		return nil
	}

	err = vsm.vaultClient.Sys().EnableAuthWithOptions(vsm.vaultConfig.KubernetesMount,
		&vaultapi.EnableAuthOptions{
			Type: model.VaultAuthTypeKubernetes.String(),
		})
	if err != nil {
		return err
	}

	return nil
}

func (vsm *VaultSystemManager) ConfigureKubernetesAuth(ctx context.Context) error {
	caData, err := vsm.GetKubernetesAuthConfig(ctx)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		model.VaultKubernetesHost:   model.DefaultVaultKubernetesHost,
		model.VaultKubernetesCACert: caData.CACertificate,
		model.VaultTokenReviewerJWT: caData.Token,
	}

	if _, err := vsm.vaultClient.Logical().Write(
		path.Join(vsm.vaultConfig.KubernetesAuthPath, model.VaultPathConfig), payload); err != nil {
		return err
	}

	return nil
}

func (vsm *VaultSystemManager) IsKubernetesAuthEnabled() (bool, error) {
	return vsm.isAuthEnabled(vsm.vaultConfig.KubernetesMount)
}

func (vsm *VaultSystemManager) GetKubernetesAuthConfig(ctx context.Context) (*model.KubernetesCertificateData, error) {
	key := types.NamespacedName{
		Name:      vsm.vaultConfig.ServiceAccount,
		Namespace: vsm.vaultConfig.Namespace,
	}

	vsa := &corev1.ServiceAccount{}
	err := vsm.kubeClient.Get(ctx, key, vsa)
	if err != nil {
		return nil, err
	}

	vaultSecret := &corev1.Secret{}
	vaultSecretKey := types.NamespacedName{
		Name:      vsa.Secrets[0].Name,
		Namespace: key.Namespace,
	}

	err = vsm.kubeClient.Get(ctx, vaultSecretKey, vaultSecret)
	if err != nil {
		return nil, err
	}

	ca := string(vaultSecret.Data[model.KubernetesSecretDataCACertificate])
	token := string(vaultSecret.Data[model.KubernetesSecretDataToken])

	return &model.KubernetesCertificateData{
		CACertificate: ca,
		Token:         token,
	}, nil
}

func (vsm *VaultSystemManager) mapSecretEngineType(vse *model.VaultSecretEngine) map[string]interface{} {
	return map[string]interface{}{
		"type": vse.Type.String(),
	}

}

func (vsm *VaultSystemManager) isAuthEnabled(mount string) (bool, error) {
	auths, err := vsm.vaultClient.Logical().Read(model.VaultSysAuth)
	if err != nil {
		return false, err
	}

	if auth, ok := auths.Data[fmt.Sprintf("%s/", mount)]; auth != nil && ok {
		return true, nil
	}

	return false, nil
}

func NewVaultSystemManager(vaultClient *vaultapi.Client, kubeClient client.Client, vaultConfig *VaultConfig) *VaultSystemManager {
	return &VaultSystemManager{
		kubeClient:  kubeClient,
		vaultClient: vaultClient,
		vaultConfig: vaultConfig,
	}
}
