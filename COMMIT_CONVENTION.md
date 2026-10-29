# Commit and PR Title Convention

Commits and PR titles follow [Conventional Commits](https://www.conventionalcommits.org).
PRs are squash-merged, so the PR title becomes the commit title on the target branch, and
the release notes are generated from these titles with [git-cliff](https://git-cliff.org)
(see [cliff.toml](cliff.toml)).

```
type(scope): description
```

- `type`: `fix`, `feat`, `perf`, `revert`, `chore`, `test`, `build`, `ci`, `docs`, `style` or `refactor`
- `scope`: optional, the area of the change, e.g. a service name like `graph` or `web`
- `!` before the colon marks a breaking change: `feat(proxy)!: remove the PROXY_LEGACY_AUTH option`
- `description`: lowercase, imperative, no trailing period

CI markers such as `[full-ci]` in front and a ticket key after the colon are allowed, they
are removed from the release notes:

```
[full-ci] fix(graph): [OCISDEV-1506] apply $filter before resolving users from other instances
```

## Release notes

Only `fix`, `feat`, `perf` and `revert` are listed, breaking changes additionally in their
own section. Each entry is the description, so it has to tell an admin what changed without
opening the PR:

| Instead of | Write |
|---|---|
| `fix: share listing` | `fix(graph): return shares of disabled spaces in the shared-with-me list` |
| `chore: fix public share` | `fix(frontend): resolve password protected public links` |
| `feat: harden nats` | `feat: add TLS support for the NATS store and registry connections` |

## Backports

Cherry-pick the change to the stable branch and keep the original PR title, so the fix
reads the same in both releases.

## Preview the release notes

Install [git-cliff](https://git-cliff.org/docs/installation), then:

```bash
make changelog-generate CHANGELOG_VERSION=v9.0.0
```
