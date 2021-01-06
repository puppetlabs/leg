// Package k8sext provides adapters for the Kubernetes API machinery clock
// implementation to this library.
package k8sext

import (
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/clock"
	testclock "k8s.io/apimachinery/pkg/util/clock"
)

type k8sClock struct{ testclock.Clock }

func (kc *k8sClock) NewTimer(d time.Duration) clock.Timer   { return kc.Clock.NewTimer(d) }
func (kc *k8sClock) NewTicker(d time.Duration) clock.Ticker { return kc.Clock.NewTicker(d) }

func NewClock(delegate testclock.Clock) clock.Clock {
	return &k8sClock{delegate}
}
