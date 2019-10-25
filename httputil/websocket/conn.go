package websocket

import (
	"io"
	"time"
)

type Conn interface {
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error

	SetPongHandler(h func(appData string) error)

	WriteControl(messageType int, data []byte, deadline time.Time) error
	NextWriter(messageType int) (io.WriteCloser, error)
	NextReader() (int, io.Reader, error)

	Close() error
}
