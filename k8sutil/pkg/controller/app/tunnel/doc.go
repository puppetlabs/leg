/*
Package tunnel provides a bidirectional HTTP stream between a cluster service
and a locally-accessible web server.

This package is useful for testing. For example, you may want to start a server
using net/http/httptest and make that server available to a test Kubernetes
cluster. To do so, simply create a tunnel:

  tun, err := ApplyHTTP(ctx, cl, key)

Then connect to it using the server URL.

  srv := httptest.NewServer(handler)
  err := WithHTTPConnection(ctx, cfg, tun, srv.URL, func(ctx context.Context) {
    // The server is now available in your cluster at tun.URL().
  })

If you need to wrap the server with TLS, you can use ApplyHTTPS instead, which
combines the HTTP tunnel with a TLS reverse proxy from the
controller/app/tlsproxy package.
*/
package tunnel
