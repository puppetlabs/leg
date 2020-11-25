package websocket_test

import (
	"testing"

	gws "github.com/gorilla/websocket"
	"github.com/puppetlabs/leg/httputil/websocket"
)

func TestConnCompatibility(t *testing.T) {
	// This test will not compile if Gorilla ever changes the methods we expect
	// on their struct.
	var _ websocket.Conn = &gws.Conn{}
}
