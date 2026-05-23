# Contributing to pipelock

Thank you for your interest in contributing to pipelock! This document outlines the process for contributing to this project.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git

### Setting Up Your Development Environment

1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/pipelock.git
   cd pipelock
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/luckyPipewrench/pipelock.git
   ```

## Making Changes

### Branching

Create a new branch for your changes:
```bash
git checkout -b feat/your-feature-name
```

Use the following prefixes:
- `feat/` — new features
- `fix/` — bug fixes
- `docs/` — documentation changes
- `refactor/` — code refactoring
- `test/` — adding or updating tests

### Code Style

- Follow standard Go formatting (`gofmt`).
- Run `go vet ./...` before submitting.
- Ensure all exported types, functions, and methods have doc comments.
- Keep functions focused and concise.

### Testing

All changes must include appropriate tests. Run the test suite with:
```bash
go test -v -race ./...
```

For benchmarks:
```bash
go test -bench=. -benchmem ./...
```

Ensure that race conditions are not introduced — the `-race` flag is mandatory for CI.

<!-- Personal note: I've been running benchmarks with -count=3 locally to get more stable results -->

### Fuzzing

This project uses ClusterFuzzLite for continuous fuzzing. If you add new parsing or input-handling logic, consider adding a corresponding fuzz target under `.clusterfuzzlite/`.

Run fuzz tests locally:
```bash
go test -fuzz=FuzzYourTarget -fuzztime=60s ./...
```

## Submitting a Pull Request

1. Ensure your branch is up to date with upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```
2. Push your branch to your fork:
   ```bash
   git push origin feat/your-feature-name
   ```
3. Open a Pull Request against the `main` branch of this repository.
4. Fill in the PR template with a clear description of your changes and the motivation behind them.
5. Link any related issues using GitHub keywords (e.g., `Closes #42`).

## Reporting Issues

Please use the GitHub Issue templates provided:
- **Bug Report** — for unexpected behavior or panics.
- **Feature Request** — for new functionality or improvements.

Include as much detail as possible, including Go version, OS, and a minimal reproducible example where applicable.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you agree to uphold these standards.

## License

By contributing, you agree that your contributions will be licensed under the same license as this project.
