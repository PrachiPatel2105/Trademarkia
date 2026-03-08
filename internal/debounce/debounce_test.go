package debounce

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestDebouncer_SingleEvent(t *testing.T) {
	var callCount int32
	callback := func() {
		atomic.AddInt32(&callCount, 1)
	}

	debouncer := New(100*time.Millisecond, callback)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go debouncer.Start(ctx)

	// Trigger once
	debouncer.Trigger()

	// Wait for debounce period + buffer
	time.Sleep(200 * time.Millisecond)

	count := atomic.LoadInt32(&callCount)
	if count != 1 {
		t.Errorf("expected 1 callback, got %d", count)
	}
}

func TestDebouncer_MultipleRapidEvents(t *testing.T) {
	var callCount int32
	callback := func() {
		atomic.AddInt32(&callCount, 1)
	}

	debouncer := New(100*time.Millisecond, callback)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go debouncer.Start(ctx)

	// Trigger multiple times rapidly
	for i := 0; i < 10; i++ {
		debouncer.Trigger()
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for debounce period + buffer
	time.Sleep(200 * time.Millisecond)

	count := atomic.LoadInt32(&callCount)
	if count != 1 {
		t.Errorf("expected 1 callback for rapid events, got %d", count)
	}
}

func TestDebouncer_SeparatedEvents(t *testing.T) {
	var callCount int32
	callback := func() {
		atomic.AddInt32(&callCount, 1)
	}

	debouncer := New(50*time.Millisecond, callback)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go debouncer.Start(ctx)

	// First event
	debouncer.Trigger()
	time.Sleep(100 * time.Millisecond)

	// Second event (after debounce period)
	debouncer.Trigger()
	time.Sleep(100 * time.Millisecond)

	count := atomic.LoadInt32(&callCount)
	if count != 2 {
		t.Errorf("expected 2 callbacks for separated events, got %d", count)
	}
}

func TestDebouncer_ContextCancellation(t *testing.T) {
	var callCount int32
	callback := func() {
		atomic.AddInt32(&callCount, 1)
	}

	debouncer := New(100*time.Millisecond, callback)
	ctx, cancel := context.WithCancel(context.Background())

	go debouncer.Start(ctx)

	// Trigger event
	debouncer.Trigger()

	// Cancel context before debounce period
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait to ensure callback doesn't fire
	time.Sleep(100 * time.Millisecond)

	count := atomic.LoadInt32(&callCount)
	if count != 0 {
		t.Errorf("expected 0 callbacks after cancellation, got %d", count)
	}
}

func TestDebouncer_HighFrequency(t *testing.T) {
	var callCount int32
	callback := func() {
		atomic.AddInt32(&callCount, 1)
	}

	debouncer := New(50*time.Millisecond, callback)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go debouncer.Start(ctx)

	// Simulate very high frequency events (like editor saves)
	for i := 0; i < 100; i++ {
		debouncer.Trigger()
		time.Sleep(1 * time.Millisecond)
	}

	// Wait for debounce
	time.Sleep(100 * time.Millisecond)

	count := atomic.LoadInt32(&callCount)
	if count != 1 {
		t.Errorf("expected 1 callback for high frequency events, got %d", count)
	}
}

func TestDebouncer_DelayAccuracy(t *testing.T) {
	var callTime time.Time
	callback := func() {
		callTime = time.Now()
	}

	delay := 100 * time.Millisecond
	debouncer := New(delay, callback)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go debouncer.Start(ctx)

	startTime := time.Now()
	debouncer.Trigger()

	// Wait for callback
	time.Sleep(200 * time.Millisecond)

	elapsed := callTime.Sub(startTime)
	
	// Allow 20ms tolerance
	if elapsed < delay || elapsed > delay+20*time.Millisecond {
		t.Errorf("expected delay ~%v, got %v", delay, elapsed)
	}
}
