# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- **New `rejonson/v9/doc` package**: typed RedisJSON documents built on go-redis v9 (`Store[T]` with generics, JSONPath helpers, NX/XX/TTL options, `Update` with optimistic locking, typed pipeline futures, `Scan` iterator, pluggable codec). Works with Client, ClusterClient, Ring and failover clients, over RESP2 and RESP3.
- Deprecated the v9 raw helpers (`JsonGet`, …, `ExtendClient`, `ExtendPipeline`) in favor of `doc` and go-redis's native JSON commands; generated via the new `Deprecated` generator flag
- README now leads with v9 and `doc`; go-redis v6–v8 modules are documented as maintenance-only
- CI tests Go 1.24.x – 1.27.x with `GOTOOLCHAIN=local`; the v9 module (go-redis v9.23, requires Go 1.26) is skipped on older toolchains
- CI runs against Redis 8 (JSON built in) and Redis Stack 7.4, replacing the `redis-stack-server:latest` image
- `test.sh` now also runs `go vet` and the race detector
- Bumped GitHub Actions (checkout v5, setup-go v6, create-pull-request v7); Dependabot now tracks Actions too
- Removed the obsolete Travis CI config
- Modernized build tooling and CI
- Added community health files (Code of Conduct, Contributing, Security policy)
- Added GitHub Actions workflow covering Go 1.20.x – 1.23.x
- Enhanced Dependabot auto-merge workflow with auto-approval for patch and minor version updates
