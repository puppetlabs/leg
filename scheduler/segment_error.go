package scheduler

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/puppetlabs/insights-stdlib/request"
	"github.com/puppetlabs/insights-stdlib/scheduler/errors"
)

var (
	// SegmentErrorBehaviorCollect allows all processes to complete and returns
	// a full set of errors when they have finished.
	SegmentErrorBehaviorCollect SegmentErrorBehavior = &segmentErrorBehaviorCollect{}

	// SegmentErrorBehaviorTerminate causes the entire segment to clean up and
	// exit when the first error occurs.
	SegmentErrorBehaviorTerminate SegmentErrorBehavior = &segmentErrorBehaviorTerminate{}

	// SegmentErrorBehaviorDrop ignores errors, merely logging them for
	// reference.
	SegmentErrorBehaviorDrop SegmentErrorBehavior = &segmentErrorBehaviorDrop{}
)

type segmentErrorHandler interface {
	ForProcess(req *request.Request, p Process, err errors.Error)
	ForDescriptor(i int, desc Descriptor, err errors.Error)
	Ch() <-chan struct{}
	Errors() []errors.Error
}

type SegmentErrorBehavior interface {
	handler() segmentErrorHandler
}

func processError(req *request.Request, p Process, err errors.Error) errors.Error {
	return errors.NewLifecycleProcessError(req.Identifier, p.Description()).WithCause(err)
}

func descriptorError(i int, desc Descriptor, err errors.Error) errors.Error {
	return errors.NewLifecycleDescriptorError(int64(i), fmt.Sprintf("%v", reflect.TypeOf(desc))).WithCause(err)
}

type segmentErrorHandlerCollect struct {
	ch   chan struct{}
	errs []errors.Error
	mut  sync.RWMutex
}

func (h *segmentErrorHandlerCollect) put(err errors.Error) {
	h.mut.Lock()
	defer h.mut.Unlock()

	h.errs = append(h.errs, err)
}

func (h *segmentErrorHandlerCollect) ForProcess(req *request.Request, p Process, err errors.Error) {
	h.put(processError(req, p, err))
}

func (h *segmentErrorHandlerCollect) ForDescriptor(i int, desc Descriptor, err errors.Error) {
	h.put(descriptorError(i, desc, err))
}

func (h *segmentErrorHandlerCollect) Ch() <-chan struct{} {
	// This handler never terminates.
	return h.ch
}

func (h *segmentErrorHandlerCollect) Errors() []errors.Error {
	h.mut.RLock()
	defer h.mut.RUnlock()

	return append([]errors.Error{}, h.errs...)
}

type segmentErrorBehaviorCollect struct{}

func (segmentErrorBehaviorCollect) handler() segmentErrorHandler {
	return &segmentErrorHandlerCollect{
		ch: make(chan struct{}),
	}
}

type segmentErrorHandlerTerminate struct {
	ch  chan struct{}
	err errors.Error
	mut sync.RWMutex
}

func (h *segmentErrorHandlerTerminate) put(err errors.Error) {
	h.mut.Lock()
	defer h.mut.Unlock()

	if h.err != nil {
		return
	}

	close(h.ch)
	h.err = err
}

func (h *segmentErrorHandlerTerminate) ForProcess(req *request.Request, p Process, err errors.Error) {
	h.put(processError(req, p, err))
}

func (h *segmentErrorHandlerTerminate) ForDescriptor(i int, desc Descriptor, err errors.Error) {
	h.put(descriptorError(i, desc, err))
}

func (h *segmentErrorHandlerTerminate) Ch() <-chan struct{} {
	return h.ch
}

func (h *segmentErrorHandlerTerminate) Errors() []errors.Error {
	h.mut.RLock()
	defer h.mut.RUnlock()

	if h.err == nil {
		return nil
	}

	return []errors.Error{h.err}
}

type segmentErrorBehaviorTerminate struct{}

func (segmentErrorBehaviorTerminate) handler() segmentErrorHandler {
	return &segmentErrorHandlerTerminate{
		ch: make(chan struct{}),
	}
}

type segmentErrorHandlerDrop struct {
	ch chan struct{}
}

func (h *segmentErrorHandlerDrop) ForProcess(req *request.Request, p Process, err errors.Error) {}
func (h *segmentErrorHandlerDrop) ForDescriptor(i int, desc Descriptor, err errors.Error)       {}
func (h *segmentErrorHandlerDrop) Ch() <-chan struct{}                                          { return h.ch }
func (h *segmentErrorHandlerDrop) Errors() []errors.Error                                       { return nil }

type segmentErrorBehaviorDrop struct{}

func (segmentErrorBehaviorDrop) handler() segmentErrorHandler {
	return &segmentErrorHandlerDrop{
		ch: make(chan struct{}),
	}
}
