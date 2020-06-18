package activities

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/puppetlabs/horsehead/v2/instrumentation/activities/activity"
)

func TestReporter(t *testing.T) {
	t.Run("with a single delegate", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		reporter := NewReporter()
		defer reporter.Close()

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.Report(NewActivity("abc123", "user-created"))

		select {
		case <-WaitGroupNotify(&wg):
			// noop
		case <-time.After(50 * time.Millisecond):
			t.Failed()
		}
	})

	t.Run("with a multiple delegates", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		reporter := NewReporter()
		defer reporter.Close()

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.Report(NewActivity("abc123", "user-created"))

		select {
		case <-WaitGroupNotify(&wg):
			// noop
		case <-time.After(50 * time.Millisecond):
			t.Failed()
		}
	})

	t.Run("with a failing delegate", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		reporter := NewReporter()
		defer reporter.Close()

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			return errors.New("oopsies")
		}))

		reporter.Report(NewActivity("abc123", "user-created"))

		select {
		case <-WaitGroupNotify(&wg):
			// noop
		case <-time.After(50 * time.Millisecond):
			t.Failed()
		}
	})
}

type mockDelegate struct {
	fn func(activity.Activity) error
}

func (d mockDelegate) Report(act activity.Activity) error {
	if d.fn != nil {
		return d.fn(act)
	}

	return nil
}

func (d mockDelegate) Close() error {
	return nil
}

func DelegateFunc(fn func(activity.Activity) error) Delegate {
	return mockDelegate{fn}
}

func WaitGroupNotify(wg *sync.WaitGroup) <-chan bool {
	ch := make(chan bool)

	go func() {
		wg.Wait()
		ch <- true
	}()

	return ch
}
