# Contributing to Rejonson

First off, thanks for taking the time to contribute! 🎉

The following is a set of guidelines for contributing to the project. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Pull Request Checklist](#pull-request-checklist)
- [Release Process](#release-process)

## Code of Conduct

This project and everyone participating in it is governed by the [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository and clone your fork locally.
2. Configure the upstream remote:
   ```bash
   git remote add upstream https://github.com/KromDaniel/rejonson.git
   ```
3. Create a feature branch from the latest `master` branch before you start working.

## Development Workflow

- Keep changes focused and minimal. If you discover unrelated issues, open an issue or a separate pull request.
- Write clear commit messages and document **why** a change is needed.
- Update documentation and examples when behavior changes.

## Testing

Tests require a Redis server with the RedisJSON module. The easiest way to set this up is via the Redis Stack Docker image:

```bash
# Start RedisJSON locally
docker run --rm -p 6379:6379 redis/redis-stack-server:latest
```

Then run the project test suite:

```bash
./test.sh
```

If you are working inside one of the submodules (`v7`, `v8`, `v9`, or `generator`), you can also run `go test ./...` within that directory.

## Pull Request Checklist

Before submitting a pull request, please ensure that:

- [ ] Tests pass locally on the supported Go versions (see the README for the list).
- [ ] New or updated APIs include clear documentation and tests.
- [ ] Public API changes preserve backwards compatibility.
- [ ] You added or updated any relevant documentation, including README examples.

## Release Process

1. Ensure `master` is green in CI.
2. Update CHANGELOG and version tags if applicable.
3. Tag the release (`git tag vX.Y.Z && git push origin vX.Y.Z`).

If you have questions, please open an issue. Happy hacking! 💻
