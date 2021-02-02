package helper

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Own(target client.Object, owner lifecycle.TypedObject) error {
	if target.GetNamespace() != owner.Object.GetNamespace() {
		return &OwnerInOtherNamespaceError{
			Owner:  owner,
			Target: target,
		}
	}

	ref := metav1.NewControllerRef(owner.Object, owner.GVK)

	targetOwners := target.GetOwnerReferences()
	for i, c := range targetOwners {
		if equality.Semantic.DeepEqual(c, *ref) {
			return nil
		} else if c.Controller != nil && *c.Controller {
			c.Controller = func(b bool) *bool { return &b }(false)
			klog.Warningf(
				"%T %s is stealing controller for %T %s from %s %s/%s",
				owner.Object, client.ObjectKeyFromObject(owner.Object),
				target, client.ObjectKeyFromObject(target),
				c.Kind, target.GetNamespace(), c.Name,
			)

			targetOwners[i] = c
		}
	}

	targetOwners = append(targetOwners, *ref)
	target.SetOwnerReferences(targetOwners)

	return nil
}

func OwnUncontrolled(target client.Object, owner lifecycle.TypedObject) error {
	if target.GetNamespace() != owner.Object.GetNamespace() {
		return &OwnerInOtherNamespaceError{
			Owner:  owner,
			Target: target,
		}
	}

	ref := metav1.OwnerReference{
		APIVersion: owner.GVK.GroupVersion().String(),
		Kind:       owner.GVK.Kind,
		Name:       owner.Object.GetName(),
		UID:        owner.Object.GetUID(),
	}

	targetOwners := target.GetOwnerReferences()
	for _, c := range targetOwners {
		if equality.Semantic.DeepEqual(c, ref) {
			return nil
		}
	}

	targetOwners = append(targetOwners, ref)
	target.SetOwnerReferences(targetOwners)

	return nil
}
