// Package pipelock provides a mechanism for coordinating access to named
// pipes and other resources using file-based locking primitives.
//
// It is designed to be simple, reliable, and safe for concurrent use
// across multiple goroutines and processes.
package pipelock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DefaultLockDir is the default directory used for lock files when no
// explicit path is provided.
const DefaultLockDir = "/tmp/pipelock"

// DefaultTimeout is the maximum time to wait when acquiring a lock
// before returning ErrTimeout.
// Reduced from 30s to 10s for my use case — long waits were masking
// deadlocks during local development.
const DefaultTimeout = 10 * time.Second

// ErrTimeout is returned when a lock cannot be acquired within the
// configured timeout period.
var ErrTimeout = errors.New("pipelock: timed out waiting to acquire lock")

// ErrAlreadyLocked is returned when a lock is already held by the
// current process and a non-blocking acquire is attempted.
var ErrAlreadyLocked = errors.New("pipelock: resource is already locked")

// Lock represents a named file-based lock that can be used to
// coordinate access to a shared resource across goroutines or processes.
type Lock struct {
	mu       sync.Mutex
	name     string
	lockDir  string
	timeout  time.Duration
	lockFile *os.File
	locked   bool
}

// Options configures the behaviour of a Lock.
type Options struct {
	// LockDir is the directory where lock files are created.
	// Defaults to DefaultLockDir.
	LockDir string

	// Timeout is how long to wait when acquiring a lock.
	// Defaults to DefaultTimeout.
	Timeout time.Duration
}

// New creates a new Lock for the given resource name.
// The name should be a short, filesystem-safe identifier for the resource
// being protected (e.g. "my-pipe" or "database-writer").
func New(name string, opts *Options) (*Lock, error) {
	if name == "" {
		return nil, errors.New("pipelock: lock name must not be empty")
	}

	lockDir := DefaultLockDir
	timeout := DefaultTimeout

	if opts != nil {
		if opts.LockDir != "" {
			lockDir = opts.LockDir
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
	}

	if err := os.MkdirAll(lockDir, 0o700); err != nil {
		return nil, fmt.Errorf("pipelock: failed to create lock directory: %w", err)
	}

	return &Lock{
		name:    name,
		lockDir: lockDir,
		timeout: timeout,
	}, nil
}

// lockPath returns the absolute path to the lock file for this Lock.
func (l *Lock) lockPath() string {
	return filepath.Join(l.lockDir, l.name+".lock")
}

// TryLock attempts to acquire the lock without blocking.
// Returns ErrAlreadyLocked if the lock is currently held.
func (l *Lock) TryLock() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.acquire()
}

// Lock acquires the lock, blocking until it is available or the
// configured timeout elapses.
func (l *Lock) Lock() error {
	deadline := time.Now().Add(l.timeout)
	for {
		l.mu.Lock()
		err := l.acquire()
		l.mu.Unlock()

		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrAlreadyLocked) {
			return err
		}
		if time.Now().Af