package storage

// Look here for various implementations:
// 	https://github.com/puppetlabs/nebula-libs/tree/master/storage"
//

import (
	"context"
	"fmt"
	"io"
)

type ErrorCode string

const (
	AuthError     ErrorCode = "AuthError"
	NotFoundError ErrorCode = "NotFoundError"
	TimeoutError  ErrorCode = "TimeoutError"
	UnknownError  ErrorCode = "UnknownError"
)

type errorImpl struct {
	message string
	code    ErrorCode
	cause   error
}

func (e *errorImpl) Error() string {
	return e.message
}

func (e *errorImpl) Unwrap() error {
	return e.cause
}

func Errorf(cause error, code ErrorCode, format string, a ...interface{}) error {
	return &errorImpl{
		code:    code,
		message: fmt.Sprintf(format, a...),
		cause:   cause,
	}
}

func IsAuthError(err error) bool {
	e, ok := err.(*errorImpl)
	return ok && e.code == AuthError
}

func IsNotFoundError(err error) bool {
	e, ok := err.(*errorImpl)
	return ok && e.code == NotFoundError
}

func IsTimeoutError(err error) bool {
	e, ok := err.(*errorImpl)
	return ok && e.code == TimeoutError
}

type Sink func(io.Writer) error
type Source func(*Meta, io.Reader) error

type Meta struct {
	ContentType string
}
type PutOptions struct {
	ContentType string
}
type GetOptions struct {
	// If Offset and Length are 0, the full blob is returned.
	//
	// If length is < 0, all bytes after Offset are returned.
	//
	// If Offset is < 0, it is the offset from the end of the blob
	// and length must be 0 or negative.
	Offset, Length int64
}
type DeleteOptions struct {
	// TODO: Support conditional deletes?
}

type BlobStore interface {
	Put(ctx context.Context, key string, sink Sink, opts PutOptions) error
	Get(ctx context.Context, key string, source Source, opts GetOptions) error
	Delete(ctx context.Context, key string, opts DeleteOptions) error
}
