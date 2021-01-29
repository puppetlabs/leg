//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen rbac:roleName=webhookcert-controller paths=./... output:artifacts:config=../../../manifests/app/webhookcert/generated

package webhookcert
