package tunnel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"github.com/rancher/remotedialer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WithHTTPConnection forwards the service provided by the given HTTP tunnel to
// the target URL accessible from the caller of this function.
//
// It invokes the given callback function when the connection is established.
// When the callback completes, the connection is torn down.
func WithHTTPConnection(ctx context.Context, cfg *rest.Config, h *HTTP, targetURL string, fn func(ctx context.Context)) (err error) {
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

	// Forward port to get access to remote side.
	req := kc.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(h.Pod.Key.Namespace).
		Name(h.Pod.Key.Name).
		SubResource("portforward").
		Param("ports", "8080")

	transport, upgrader, err := spdy.RoundTripperFor(cfg)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())
	readyCh := make(chan struct{})

	pf, err := portforward.New(dialer, []string{":8080"}, ctx.Done(), readyCh, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}
	defer pf.Close()

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

	// Connect local side.
	ports, err := pf.GetPorts()
	if err != nil {
		return err
	} else if len(ports) != 1 {
		return errors.New("missing local port information")
	}

	go func() {
		headers := make(http.Header)
		headers.Set("x-inlets-id", uuid.New().String())
		headers.Set("x-inlets-upstream", "="+targetURL)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			remotedialer.ClientConnect(
				ctx,
				fmt.Sprintf("ws://localhost:%d/tunnel", ports[0].Local),
				headers,
				nil,
				func(proto, address string) bool { return true },
				nil,
			)
		}
	}()

	// Wait for service.
	if _, err := corev1obj.NewEndpointsBoundPoller(corev1obj.NewEndpoints(h.Service)).Load(ctx, cl); err != nil {
		return err
	}

	fn(ctx)
	return
}
