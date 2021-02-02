//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen rbac:roleName=selfsignedsecret-controller paths=./... output:artifacts:config=../../../manifests/app/selfsignedsecret/generated

package selfsignedsecret
