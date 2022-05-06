package tunnel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/puppetlabs/leg/k8sutil/pkg/app/portforward"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/rancher/remotedialer"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WithHTTPConnection forwards the service provided by the given HTTP tunnel to
// the target URL accessible from the caller of this function.
//
// It invokes the given callback function when the connection is established.
// When the callback completes, the connection is torn down.
func WithHTTPConnection(ctx context.Context, cfg *rest.Config, h *HTTP, targetURL string, fn func(ctx context.Context)) error {
	cl, err := client.New(cfg, client.Options{})
	if err != nil {
		return err
	}

	if _, err := corev1obj.NewEndpointsBoundPoller(corev1obj.NewEndpoints(h.Service)).Load(ctx, cl); err != nil {
		return err
	}

	return portforward.ForwardPod(ctx, cfg, h.Pod, []uint16{8080}, func(ctx context.Context, m portforward.Map) error {
		connCh := make(chan struct{})

		go func(connCh chan struct{}) {
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
					fmt.Sprintf("ws://localhost:%d/tunnel", m[8080]),
					headers,
					nil,
					func(proto, address string) bool { return true },
					func(ctx context.Context) error {
						if connCh != nil {
							close(connCh)
							connCh = nil
						}
						return nil
					},
				)
			}
		}(connCh)

		// Wait for client connection.
		select {
		case <-connCh:
		case <-ctx.Done():
			return ctx.Err()
		}

		fn(ctx)
		return nil
	})
}
