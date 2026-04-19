# agents.md -- ownCloud Infinite Scale

## Repository Overview

The core oCIS repository -- a next-generation file sync and share platform written in Go. Single-binary, cloud-native microservices architecture. Licensed under Apache-2.0.

## Architecture & Key Paths

- `ocis/` -- Main binary entry point
- `services/` -- Individual microservice implementations
- `ocis-pkg/` -- Shared packages across services
- `internal/` -- Internal packages
- `protogen/` -- Generated protobuf code
- `vendor/` -- Vendored Go dependencies
- `deployments/` -- Docker Compose and deployment configs
- `tests/` -- Acceptance and integration tests
- `docs/` -- Documentation source
- `scripts/` -- Build and utility scripts
- `tools/` -- Development tools
- `assets/` -- Static assets
- `Makefile` -- Build and test automation
- `Dockerfile` -- Docker image build
- `go.mod` / `go.sum` -- Go module definition
- `composer.json` -- PHP dependencies (for acceptance tests)
- `vendor-bin/` -- PHP tools (behat)

## Development Conventions

- Go codebase with microservices architecture
- Protobuf for service definitions
- Makefile-driven build system
- Acceptance tests use Behat (PHP)
- SonarCloud for quality gate

## Build & Test Commands

```bash
make -C ocis build            # Build the ocis binary
make test                     # Run Go tests
make go-coverage              # Generate coverage report
make generate                 # Run code generation
make vet                      # Run go vet
make ci-go-generate           # CI code generation
make go-mod-tidy              # Tidy go modules
```

## Important Constraints

- Licensed under Apache-2.0 (already at the OSPO target license). The broader ownCloud organization is migrating other repositories from copyleft licenses to Apache 2.0.
- Contains protobuf-generated code in `protogen/` -- do not edit directly.
- All contributions require a DCO sign-off.


## OSPO Policy Constraints

### GitHub Actions
- **Only** use actions owned by `owncloud`, created by GitHub (`actions/*`), or verified on the GitHub Marketplace.
- Pin all actions to their full commit SHA (not tags): `uses: actions/checkout@<SHA> # vX.Y.Z`
- Never introduce actions from unverified third parties.

### Dependency Management
- Dependabot is configured for automated dependency updates.
- Review and merge Dependabot PRs as part of regular maintenance.
- Do not introduce new dependencies without discussion in an issue first.

### Git Workflow
- **Rebase policy**: Always rebase; never create merge commits. Use `git pull --rebase` and `git rebase` before pushing.
- **Signed commits**: All commits **must** be PGP/GPG signed (`git commit -S -s`).
- **DCO sign-off**: Every commit needs a `Signed-off-by` line (`git commit -s`).
- **Conventional Commits**: Use the [Conventional Commits](https://www.conventionalcommits.org/) format where the repository enforces it.

## Context for AI Agents

oCIS is the primary ownCloud product. It uses a microservices architecture where each service in `services/` runs as part of a single binary. The `ocis-pkg/` directory contains shared libraries. Protobuf definitions generate service interfaces. The Libre Graph API is the primary API surface.
