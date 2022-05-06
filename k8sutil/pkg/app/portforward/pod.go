package portforward

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Map map[uint16]uint16

// ForwardPod connects to the given pod on an HTTP-accessible Kubernetes
// instance, exposing the specified remote TCP ports on the pod locally.
//
// It invokes the given callback function with a map of the remote port number
// to the allocated local port number. When the callback completes, the
// connection is torn down.
func ForwardPod(ctx context.Context, cfg *rest.Config, pod *corev1obj.Pod, ports []uint16, fn func(ctx context.Context, m Map) error) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cl, err := client.New(cfg, client.Options{})
	if err != nil {
		return err
	}

	kc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	_, err = corev1obj.NewPodRunningPoller(pod).Load(ctx, cl)
	if err != nil {
		return err
	}

	req := kc.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(pod.Key.Namespace).
		Name(pod.Key.Name).
		SubResource("portforward")
	for _, port := range ports {
		req = req.Param("ports", strconv.FormatUint(uint64(port), 10))
	}

	transport, upgrader, err := spdy.RoundTripperFor(cfg)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())
	readyCh := make(chan struct{})

	fwds := make([]string, len(ports))
	for i, port := range ports {
		fwds[i] = fmt.Sprintf(":%d", port)
	}
	pf, err := portforward.New(dialer, fwds, ctx.Done(), readyCh, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}

	fwdCh := retry.WaitAsync(ctx, func(ctx context.Context) (bool, error) {
		err := pf.ForwardPorts()
		return err == nil, err
	})
	defer func() {
		select {
		case ferr := <-fwdCh:
			if err == nil {
				err = ferr
			}
		default:
		}
	}()

	select {
	case <-readyCh:
	case <-ctx.Done():
		return ctx.Err()
	}

	conns, err := pf.GetPorts()
	if err != nil {
		return err
	}

	m := make(Map)
	for _, conn := range conns {
		m[conn.Remote] = conn.Local
	}

	return fn(ctx, m)
}
