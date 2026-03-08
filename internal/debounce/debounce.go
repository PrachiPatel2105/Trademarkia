package debounce

import (
	"context"
	"time"
)

// Debouncer aggregates rapid events and triggers a single callback
type Debouncer struct {
	delay    time.Duration
	timer    *time.Timer
	callback func()
	eventCh  chan struct{}
}

// New creates a new Debouncer
func New(delay time.Duration, callback func()) *Debouncer {
	return &Debouncer{
		delay:    delay,
		callback: callback,
		eventCh:  make(chan struct{}, 100),
	}
}

// Trigger signals that an event occurred
func (d *Debouncer) Trigger() {
	select {
	case d.eventCh <- struct{}{}:
	default:
		// Channel full, skip this event
	}
}

// Start begins processing events
func (d *Debouncer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if d.timer != nil {
				d.timer.Stop()
			}
			return
		case <-d.eventCh:
			// Reset timer on each event
			if d.timer != nil {
				d.timer.Stop()
			}
			d.timer = time.AfterFunc(d.delay, d.callback)
		}
	}
}
