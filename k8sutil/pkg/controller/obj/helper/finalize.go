package helper

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func AddFinalizer(target metav1.Object, name string) bool {
	finalizers := target.GetFinalizers()
	for _, f := range finalizers {
		if f == name {
			return false
		}
	}

	target.SetFinalizers(append(finalizers, name))
	return true
}

func RemoveFinalizer(target metav1.Object, name string) bool {
	finalizers := target.GetFinalizers()
	cut := -1
	for i, f := range finalizers {
		if f == name {
			cut = i
			break
		}
	}

	if cut < 0 {
		return false
	}

	finalizers[cut] = finalizers[len(finalizers)-1]
	target.SetFinalizers(finalizers[:len(finalizers)-1])
	return true
}
