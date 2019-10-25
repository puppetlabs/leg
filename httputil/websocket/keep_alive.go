package websocket

import (
	"io"
	"time"

	gws "github.com/gorilla/websocket"
)

var (
	KeepAlivePeriod  = 20 * time.Second
	KeepAliveTimeout = 30 * time.Second

	WriteDeadline = 10 * time.Second
)

type kaConn struct {
	Conn

	ticker *time.Ticker
}

func (kc *kaConn) NextWriter(messageType int) (io.WriteCloser, error) {
	kc.SetWriteDeadline(time.Now().Add(WriteDeadline))
	return kc.Conn.NextWriter(messageType)
}

func (kc *kaConn) keepAlive() {
	for range kc.ticker.C {
		if err := kc.WriteControl(gws.PingMessage, []byte{}, time.Now().Add(WriteDeadline)); err != nil {
			break
		}
	}
}

func (kc *kaConn) Close() error {
	if kc.ticker != nil {
		kc.ticker.Stop()
	}

	return kc.Conn.Close()
}

func NewKeepAliveConn(conn Conn) Conn {
	if _, ok := conn.(*kaConn); ok {
		return conn
	}

	kc := &kaConn{
		Conn:   conn,
		ticker: time.NewTicker(KeepAlivePeriod),
	}

	kc.SetReadDeadline(time.Now().Add(KeepAliveTimeout))
	kc.SetPongHandler(func(string) error {
		kc.SetReadDeadline(time.Now().Add(KeepAliveTimeout))
		return nil
	})

	go kc.keepAlive()

	return kc
}
