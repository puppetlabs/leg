package portforward

import (
	"context"
	"errors"
	"fmt"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ForwardService forwards a single port from a given service on an
// HTTP-accessible Kubernetes instance. It always picks the first pod that
// serves as an endpoint for that service and port combination, and otherwise
// behaves like ForwardPod.
func ForwardService(ctx context.Context, cfg *rest.Config, svc *corev1obj.Service, port uint16, fn func(ctx context.Context, port uint16) error) error {
	targetPort, found := servicePortToPodPort(svc, port)
	if !found {
		return fmt.Errorf("port %d is not defined for service %q", port, svc.Key)
	}

	cl, err := client.New(cfg, client.Options{})
	if err != nil {
		return err
	}

	eps := corev1obj.NewEndpoints(svc)

	if _, err := corev1obj.NewEndpointsBoundPoller(eps).Load(ctx, cl); err != nil {
		return err
	}

	fwdAddrs, fwdPort := endpointForTargetPort(eps, targetPort)

	pod, err := firstPodInEndpointAddresses(fwdAddrs)
	if err != nil {
		return err
	}

	return ForwardPod(ctx, cfg, pod, []uint16{uint16(fwdPort.Port)}, func(ctx context.Context, m Map) error {
		return fn(ctx, m[uint16(fwdPort.Port)])
	})
}

func servicePortToPodPort(svc *corev1obj.Service, port uint16) (intstr.IntOrString, bool) {
	for _, candidate := range svc.Object.Spec.Ports {
		if candidate.Protocol != corev1.ProtocolTCP {
			continue
		}

		if candidate.Port == int32(port) {
			return candidate.TargetPort, true
		}
	}

	return intstr.IntOrString{}, false
}

func endpointForTargetPort(eps *corev1obj.Endpoints, port intstr.IntOrString) ([]corev1.EndpointAddress, corev1.EndpointPort) {
	for _, subset := range eps.Object.Subsets {
		for _, candidate := range subset.Ports {
			if candidate.Protocol != corev1.ProtocolTCP {
				continue
			}

			if port.Type == intstr.Int {
				if int(candidate.Port) != port.IntValue() {
					continue
				}
			} else {
				if s := port.String(); s != "" && candidate.Name != s {
					continue
				}
			}

			return subset.Addresses, candidate
		}
	}

	return nil, corev1.EndpointPort{}
}

func firstPodInEndpointAddresses(addrs []corev1.EndpointAddress) (*corev1obj.Pod, error) {
	for _, addr := range addrs {
		if addr.TargetRef == nil {
			continue
		}

		gv, err := schema.ParseGroupVersion(addr.TargetRef.APIVersion)
		if err != nil {
			return nil, err
		} else if !(gv == schema.GroupVersion{} || gv == corev1.SchemeGroupVersion) || addr.TargetRef.Kind != "Pod" {
			continue
		}

		return corev1obj.NewPod(client.ObjectKey{
			Namespace: addr.TargetRef.Namespace,
			Name:      addr.TargetRef.Name,
		}), nil
	}

	return nil, errors.New("no pods bound")
}
