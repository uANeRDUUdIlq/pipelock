package pipelock

import (
	"sync"
	"testing"
	"time"
)

// TestNew verifies that a new PipeLock is created successfully.
func TestNew(t *testing.T) {
	pl := New()
	if pl == nil {
		t.Fatal("expected non-nil PipeLock, got nil")
	}
}

// TestLockUnlock verifies basic lock and unlock functionality.
func TestLockUnlock(t *testing.T) {
	pl := New()

	pl.Lock()
	pl.Unlock()
}

// TestRLockRUnlock verifies basic read lock and read unlock functionality.
func TestRLockRUnlock(t *testing.T) {
	pl := New()

	pl.RLock()
	pl.RUnlock()
}

// TestConcurrentReaders verifies that multiple readers can hold the lock simultaneously.
func TestConcurrentReaders(t *testing.T) {
	pl := New()
	const numReaders = 10

	var wg sync.WaitGroup
	start := make(chan struct{})
	active := make(chan struct{}, numReaders)

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			pl.RLock()
			active <- struct{}{}
			time.Sleep(10 * time.Millisecond)
			<-active
			pl.RUnlock()
		}()
	}

	close(start)
	wg.Wait()
}

// TestWriterExcludesReaders verifies that a writer blocks readers.
func TestWriterExcludesReaders(t *testing.T) {
	pl := New()

	pl.Lock()

	readerDone := make(chan struct{})
	go func() {
		pl.RLock()
		close(readerDone)
		pl.RUnlock()
	}()

	// Give the goroutine time to attempt acquiring the read lock.
	select {
	case <-readerDone:
		t.Fatal("reader acquired lock while writer held it")
	case <-time.After(50 * time.Millisecond):
		// Expected: reader is blocked.
	}

	pl.Unlock()

	select {
	case <-readerDone:
		// Expected: reader acquired lock after writer released it.
	case <-time.After(500 * time.Millisecond):
		t.Fatal("reader did not acquire lock after writer released it")
	}
}

// TestReaderExcludesWriter verifies that readers block a writer.
func TestReaderExcludesWriter(t *testing.T) {
	pl := New()

	pl.RLock()

	writerDone := make(chan struct{})
	go func() {
		pl.Lock()
		close(writerDone)
		pl.Unlock()
	}()

	select {
	case <-writerDone:
		t.Fatal("writer acquired lock while reader held it")
	case <-time.After(50 * time.Millisecond):
		// Expected: writer is blocked.
	}

	pl.RUnlock()

	select {
	case <-writerDone:
		// Expected.
	case <-time.After(500 * time.Millisecond):
		t.Fatal("writer did not acquire lock after reader released it")
	}
}

// TestTryLock verifies non-blocking lock acquisition.
func TestTryLock(t *testing.T) {
	pl := New()

	if !pl.TryLock() {
		t.Fatal("expected TryLock to succeed on unlocked PipeLock")
	}

	if pl.TryLock() {
		t.Fatal("expected TryLock to fail while lock is held")
	}

	pl.Unlock()

	if !pl.TryLock() {
		t.Fatal("expected TryLock to succeed after unlock")
	}
	pl.Unlock()
}

// TestTryRLock verifies non-blocking read lock acquisition.
func TestTryRLock(t *testing.T) {
	pl := New()

	if !pl.TryRLock() {
		t.Fatal("expected TryRLock to succeed on unlocked PipeLock")
	}

	// A second reader should also succeed.
	if !pl.TryRLock() {
		t.Fatal("expected TryRLock to succeed when only readers hold the lock")
	}

	pl.RUnlock()
	pl.RUnlock()

	// Writer holds lock — TryRLock should fail.
	pl.Lock()
	if pl.TryRLock() {
		t.Fatal("expected TryRLock to fail while writer holds the lock")
	}
	pl.Unlock()
}
