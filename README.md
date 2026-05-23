# pipelock

[![Go Reference](https://pkg.go.dev/badge/github.com/luckyPipewrench/pipelock.svg)](https://pkg.go.dev/github.com/luckyPipewrench/pipelock)
[![Go Report Card](https://goreportcard.com/badge/github.com/luckyPipewrench/pipelock)](https://goreportcard.com/report/github.com/luckyPipewrench/pipelock)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, pipe-based read/write mutex for Go that provides named locking with support for concurrent readers and exclusive writers.

## Features

- Named locks — acquire locks by string key
- Reader/writer semantics — multiple concurrent readers, exclusive writers
- Deadlock-safe — non-blocking `TryLock` and `TryRLock` variants
- Zero external dependencies
- Fuzz tested via ClusterFuzzLite

## Installation

```bash
go get github.com/luckyPipewrench/pipelock
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/luckyPipewrench/pipelock"
)

func main() {
    pl := pipelock.New()

    // Exclusive write lock
    pl.Lock("my-resource")
    fmt.Println("writing...")
    pl.Unlock("my-resource")

    // Shared read lock
    pl.RLock("my-resource")
    fmt.Println("reading...")
    pl.RUnlock("my-resource")
}
```

## API

### `New() *PipeLock`

Creates and returns a new `PipeLock` instance.

### `Lock(key string)`

Acquires an exclusive write lock for the given key. Blocks until the lock is available.

### `Unlock(key string)`

Releases the exclusive write lock for the given key.

### `RLock(key string)`

Acquires a shared read lock for the given key. Multiple goroutines may hold a read lock simultaneously.

### `RUnlock(key string)`

Releases a shared read lock for the given key.

### `TryLock(key string) bool`

Attempts to acquire an exclusive write lock without blocking. Returns `true` if successful, `false` otherwise.

### `TryRLock(key string) bool`

Attempts to acquire a shared read lock without blocking. Returns `true` if successful, `false` otherwise.

## Concurrency Model

`pipelock` follows standard read/write mutex semantics:

- Any number of goroutines can hold a **read lock** on the same key simultaneously.
- Only one goroutine can hold a **write lock** on a given key at a time.
- A write lock cannot be acquired while any read locks are held, and vice versa.

> **Note (personal):** I've found it useful to pair `TryLock` with a short retry loop and
> `time.Sleep` backoff when contention is expected but a hard block is undesirable — avoids
> goroutine stacking without needing a separate timeout context.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on submitting issues and pull requests.

## License

MIT — see [LICENSE](LICENSE) for details.
