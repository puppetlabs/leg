package activities

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/puppetlabs/leg/instrumentation/activities/activity"
	"github.com/puppetlabs/leg/scheduler"
)

func TestReporter(t *testing.T) {
	t.Run("with a single delegate", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		reporter := NewReporter()

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.Scheduler().Start(scheduler.LifecycleStartOptions{})

		reporter.Report(context.Background(), NewActivity("abc123", "user-created"))

		select {
		case <-WaitGroupNotify(&wg):
			// noop
		case <-time.After(50 * time.Millisecond):
			t.FailNow()
		}
	})

	t.Run("with a multiple delegates", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		reporter := NewReporter()

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.Scheduler().Start(scheduler.LifecycleStartOptions{})

		reporter.Report(context.Background(), NewActivity("abc123", "user-created"))

		select {
		case <-WaitGroupNotify(&wg):
			// noop
		case <-time.After(50 * time.Millisecond):
			t.FailNow()
		}
	})

	t.Run("with a failing delegate", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		reporter := NewReporter()

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			wg.Done()
			return nil
		}))

		reporter.AddDelegate(DelegateFunc(func(activity.Activity) error {
			return errors.New("oopsies")
		}))

		reporter.Scheduler().Start(scheduler.LifecycleStartOptions{})

		reporter.Report(context.Background(), NewActivity("abc123", "user-created"))

		select {
		case <-WaitGroupNotify(&wg):
			// noop
		case <-time.After(50 * time.Millisecond):
			t.FailNow()
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
