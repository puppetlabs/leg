package corev1

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/pem"
	"errors"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrNotOpaqueSecret                = errors.New("secret is not an unstructured opaque secret")
	ErrNotImagePullSecret             = errors.New("secret is not usable for pulling container images")
	ErrNotServiceAccountTokenSecret   = errors.New("secret is not usable for service accounts")
	ErrServiceAccountTokenMissingData = errors.New("service account token secret has no token data")
	ErrNotTLSSecret                   = errors.New("secret is not usable for TLS")
)

var (
	SecretKind = corev1.SchemeGroupVersion.WithKind("Secret")
)

type Secret struct {
	Key    client.ObjectKey
	Object *corev1.Secret
}

var _ lifecycle.Deleter = &Secret{}
var _ lifecycle.LabelAnnotatableFrom = &Secret{}
var _ lifecycle.Loader = &Secret{}
var _ lifecycle.Ownable = &Secret{}
var _ lifecycle.Persister = &Secret{}

func (s *Secret) Delete(ctx context.Context, cl client.Client) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, s.Object)
}

func (s *Secret) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&s.Object.ObjectMeta, from)
}

func (s *Secret) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, s.Key, s.Object)
}

func (s *Secret) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(ctx, s.Object, owner)
}

func (s *Secret) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, s.Object, helper.WithObjectKey(s.Key)); err != nil {
		return err
	}

	s.Key = client.ObjectKeyFromObject(s.Object)
	return nil
}

func (s *Secret) Copy() *Secret {
	return &Secret{
		Key:    s.Key,
		Object: s.Object.DeepCopy(),
	}
}

func NewSecret(key client.ObjectKey) *Secret {
	return &Secret{
		Key:    key,
		Object: &corev1.Secret{},
	}
}

func NewSecretFromObject(obj *corev1.Secret) *Secret {
	return &Secret{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewSecretPatcher(upd, orig *Secret) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}

type OpaqueSecret struct {
	*Secret
}

func (os *OpaqueSecret) Load(ctx context.Context, cl client.Client) (bool, error) {
	ok, err := os.Secret.Load(ctx, cl)
	if err != nil {
		return false, err
	}

	if os.Object.Type != corev1.SecretTypeOpaque {
		return false, ErrNotOpaqueSecret
	}

	return ok, nil
}

func (os *OpaqueSecret) Copy() *OpaqueSecret {
	return &OpaqueSecret{
		Secret: os.Secret.Copy(),
	}
}

func (os *OpaqueSecret) Data(key string) (string, bool) {
	b, found := os.Object.Data[key]
	if !found {
		return "", false
	}

	return string(b), true
}

func NewOpaqueSecret(key client.ObjectKey) *OpaqueSecret {
	s := NewSecret(key)
	s.Object.Type = corev1.SecretTypeOpaque

	return &OpaqueSecret{
		Secret: s,
	}
}

func NewOpaqueSecretPatcher(upd, orig *OpaqueSecret) lifecycle.Persister {
	return NewSecretPatcher(upd.Secret, orig.Secret)
}

type ImagePullSecret struct {
	*Secret
}

func (ips *ImagePullSecret) Load(ctx context.Context, cl client.Client) (bool, error) {
	ok, err := ips.Secret.Load(ctx, cl)
	if err != nil {
		return false, err
	}

	if ips.Object.Type != corev1.SecretTypeDockerConfigJson {
		return false, ErrNotImagePullSecret
	}

	return ok, nil
}

func (ips *ImagePullSecret) Copy() *ImagePullSecret {
	return &ImagePullSecret{
		Secret: ips.Secret.Copy(),
	}
}

func (ips *ImagePullSecret) CopyFrom(src *ImagePullSecret) {
	ips.Object.Data = src.Copy().Object.Data
}

func NewImagePullSecret(key client.ObjectKey) *ImagePullSecret {
	s := NewSecret(key)
	s.Object.Type = corev1.SecretTypeDockerConfigJson

	return &ImagePullSecret{
		Secret: s,
	}
}

func NewImagePullSecretPatcher(upd, orig *ImagePullSecret) lifecycle.Persister {
	return NewSecretPatcher(upd.Secret, orig.Secret)
}

type ServiceAccountTokenSecret struct {
	*Secret
}

func (sats *ServiceAccountTokenSecret) Load(ctx context.Context, cl client.Client) (bool, error) {
	ok, err := sats.Secret.Load(ctx, cl)
	if err != nil {
		return false, err
	}

	if sats.Object.Type != corev1.SecretTypeServiceAccountToken {
		return false, ErrNotServiceAccountTokenSecret
	}

	return ok, nil
}

func (sats *ServiceAccountTokenSecret) Copy() *ServiceAccountTokenSecret {
	return &ServiceAccountTokenSecret{
		Secret: sats.Secret.Copy(),
	}
}

func (sats *ServiceAccountTokenSecret) Token() (string, error) {
	tok := string(sats.Object.Data["token"])
	if tok == "" {
		return "", ErrServiceAccountTokenMissingData
	}

	return tok, nil
}

func NewServiceAccountTokenSecret(key client.ObjectKey) *ServiceAccountTokenSecret {
	s := NewSecret(key)
	s.Object.Type = corev1.SecretTypeServiceAccountToken

	return &ServiceAccountTokenSecret{
		Secret: s,
	}
}

func NewServiceAccountTokenSecretPatcher(upd, orig *ServiceAccountTokenSecret) lifecycle.Persister {
	return NewSecretPatcher(upd.Secret, orig.Secret)
}

type TLSSecret struct {
	*Secret
}

func (ts *TLSSecret) Load(ctx context.Context, cl client.Client) (bool, error) {
	ok, err := ts.Secret.Load(ctx, cl)
	if err != nil {
		return false, err
	}

	if ts.Object.Type != corev1.SecretTypeTLS {
		return false, ErrNotTLSSecret
	}

	return ok, nil
}

func (ts *TLSSecret) Copy() *TLSSecret {
	return &TLSSecret{
		Secret: ts.Secret.Copy(),
	}
}

// Certificate returns the TLS certificate encoded in this secret. If the secret
// contains a ca.crt key that does not also exist in the tls.crt, it will be
// appended to the certificate bundle.
func (ts *TLSSecret) Certificate() (tls.Certificate, error) {
	certPEM := ts.Object.Data["tls.crt"]
	keyPEM := ts.Object.Data["tls.key"]

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return cert, err
	}

	caDER, _ := pem.Decode(ts.Object.Data["ca.crt"])
	if caDER != nil && caDER.Type == "CERTIFICATE" && !bytes.Equal(caDER.Bytes, cert.Certificate[len(cert.Certificate)-1]) {
		cert.Certificate = append(cert.Certificate, caDER.Bytes)
	}

	return cert, nil
}

func NewTLSSecret(key client.ObjectKey) *TLSSecret {
	s := NewSecret(key)
	s.Object.Type = corev1.SecretTypeTLS

	return &TLSSecret{
		Secret: s,
	}
}
