# Commit and PR Title Convention

The release notes on the [GitHub releases page](https://github.com/owncloud/ocis/releases)
and in [CHANGELOG.md](CHANGELOG.md) are generated from the history of `master` and the
stable branches with [git-cliff](https://git-cliff.org) (see [cliff.toml](cliff.toml)).
**Every merged PR becomes one line in the release notes, and that line is the PR title.**

## What ends up in the history

PRs are squash-merged, so each PR becomes a single commit on the target branch:

- **Commit title** = PR title, with the PR number appended by GitHub: `fix(graph): Fix share listing (#13154)`
- **Commit body** = PR description

The commits on your PR branch are not kept, so name them as you like. The **PR title** and
the **PR description** are what count, and both can still be edited before merging.

## PR title format

```
type(scope): Title
```

| Part | Required | Rules |
|---|---|---|
| `type` | yes | one of the types below, lowercase |
| `(scope)` | no | the area of the change, see [Scopes](#scopes) |
| `!` | only for breaking changes | right before the colon: `feat(proxy)!:` |
| `Title` | yes | what changed, written for someone reading the release notes |

A ticket key and CI markers may be added after the colon. They stay in the git history but
are removed from the release notes (Jira is not public):

```
fix(web): [OCISDEV-1541] Redirect vault root to the default app
feat: [full-ci] [OCISDEV-1563] Remove the owncloudsql driver
```

The `Validate PR title` check enforces the format, the type and the scope.

### Types

| Type | Use for | Release notes section |
|---|---|---|
| `fix` | a bug fix | Bug Fixes |
| `feat` | a new feature or a user-visible improvement | Features |
| `perf` | a performance improvement | Performance |
| `security` | a security fix or a security relevant dependency bump | Security |
| `revert` | reverting an earlier change | Reverts |
| `chore` | maintenance, version bumps, translations | not listed |
| `build` | build system, dependencies (dependabot) | not listed |
| `ci` | CI workflows | not listed |
| `test` | tests only | not listed |
| `docs` | documentation only | not listed |
| `refactor` | code changes without a behaviour change | not listed |
| `style` | formatting only | not listed |

**Pick the type by what the user notices, not by the kind of work.** A bug fix titled
`chore:` does not appear in the release notes at all. A refactoring titled `fix:` shows
up as a bug fix that nobody can find.

`fix` and `feat` PRs that only touch tests (title contains `test`, `e2e` or `flaky`) are
left out of the release notes, use `test` for them.

### Scopes

The scope is the area the change belongs to. It is printed in bold in front of the entry,
so readers can find the changes for the part they care about.

- a service name from [`services/`](services), e.g. `graph`, `proxy`, `storage-users`, `web`
- `ocis` for changes across several services or to the `ocis` binary itself
- `deps`, `deps-dev`, `deployments`, `ci`, `release`, `tests`

Any other scope fails the title check, including a ticket key such as `fix(OCISDEV-123):`.
Leave the scope out only when no scope fits.

### Writing the title

The title is read in the release notes, next to a hundred others. Write it for an admin
who is about to upgrade:

- start with a verb in the imperative: *Fix*, *Add*, *Remove*, *Return*, *Allow*
- say what changed for the user, not how it was done
- one change per PR, if the title needs an "and", consider splitting the PR
- no trailing period

| Instead of | Write |
|---|---|
| `fix: bug` | `fix(graph): Return 404 instead of 500 for a deleted space` |
| `fix: address review comments` | the title of what the PR actually changes |
| `feat: harden nats communication` | `feat(ocis): Add TLS support for NATS store and registry connections` |
| `chore: fix public share` | `fix(frontend): Fix public share resolution for password protected links` |
| `Fix: ocis_multi deployment` | `fix(deployments): Create LDAP groups in the ocis_multi example` |
| `fix(OCISDEV-123): Fix login` | `fix(proxy): [OCISDEV-123] Fix login`, the ticket is not a scope |

## PR description

The description is free text for the reviewers. Two kinds of lines in it are also read by
the release notes, anywhere in the description:

| Line | Effect |
|---|---|
| `Fixes #13100` | links the GitHub issue: `..., closes #13100` |
| `BREAKING CHANGE: <what changed and what to do>` | adds the text to the Breaking Changes section |

## Breaking changes

A change is breaking when an admin or API client has to do something when upgrading:
a removed or renamed environment variable, a changed default, a removed API endpoint,
a required migration.

1. Add `!` after the type or scope: `feat(proxy)!: Remove the PROXY_LEGACY_AUTH option`
2. Add a `BREAKING CHANGE:` line to the description that tells the reader what to do

Breaking changes are listed at the top of the release, because they are what admins read
first.

## Example

PR title:

```
fix(graph): [OCISDEV-1506] Enforce $filter before cross-instance user resolution
```

PR description:

```
The user listing resolved users from other instances before applying $filter,
which returned users that did not match the filter.

Fixes #13100
```

A breaking change in the same release, PR title:

```
feat(proxy)!: [OCISDEV-1600] Remove the PROXY_LEGACY_AUTH option
```

PR description ends with:

```
BREAKING CHANGE: PROXY_LEGACY_AUTH was removed. Remove it from your configuration, legacy basic auth is no longer supported.
```

Release notes:

```markdown
### Breaking Changes

* **proxy:** PROXY_LEGACY_AUTH was removed. Remove it from your configuration, legacy basic auth is no longer supported.

### Bug Fixes

* **graph:** Enforce $filter before cross-instance user resolution ([#13154](https://github.com/owncloud/ocis/pull/13154)), closes [#13100](https://github.com/owncloud/ocis/issues/13100)

### Features

* **proxy:** Remove the PROXY_LEGACY_AUTH option ([#13300](https://github.com/owncloud/ocis/pull/13300))
```

## What wrong titles look like

The release notes take the title as it is. These PRs, all merged into the same release:

| # | PR title | Problem |
|---|---|---|
| 1 | `Fix login for users with umlauts` | no type |
| 2 | `chore(proxy): Fix login for users with umlauts` | wrong type, it is a fix |
| 3 | `fix(web): address review comments` | title describes the PR, not the change |
| 4 | `fix(graph): Fix share listing and add TLS for NATS and bump reva` | several changes in one PR |
| 5 | `feat(proxy): Remove the PROXY_LEGACY_AUTH option` | breaking, but no `!` and no `BREAKING CHANGE:` |
| 6 | `fix(web): Fix flaky upload in the files app` | contains `flaky`, treated as a test fix |
| 7 | `Revert "fix(web): Fix the upload dialog"` | GitHub's default revert title |

become:

```markdown
### Bug Fixes

* **web:** Address review comments ([#3](https://github.com/owncloud/ocis/pull/3))
* **graph:** Fix share listing and add TLS for NATS and bump reva ([#4](https://github.com/owncloud/ocis/pull/4))

### Features

* **proxy:** Remove the PROXY_LEGACY_AUTH option ([#5](https://github.com/owncloud/ocis/pull/5))

### Miscellaneous

* Fix login for users with umlauts ([#1](https://github.com/owncloud/ocis/pull/1))
* Revert "fix(web): Fix the upload dialog" ([#7](https://github.com/owncloud/ocis/pull/7))
```

- **#2 and #6 are missing.** Users never learn that the fixes were released.
- **#3** tells the reader nothing.
- **#4** hides the TLS feature and the reva bump inside a bug fix line.
- **#5** is not in a Breaking Changes section, so admins upgrade without removing the option.
- **#1 and #7** land in *Miscellaneous*, the section for titles that could not be parsed.
  Retitle reverts to `revert(web): Restore the old upload dialog`.

The `Validate PR title` check rejects #1 and #7. #2 to #6 have a valid format, so only the
review catches them: **check the title like you check the code.**

## Backports

A fix is merged to `master` first and then cherry-picked to the stable branches that need it.

```bash
git fetch origin
git switch -c backport-13154-stable-8.2 origin/stable-8.2
git cherry-pick -x <sha of the squash commit of #13154 on master>
git push -u origin backport-13154-stable-8.2
```

Open the PR against `stable-8.2` with the **original title** and the branch in front:

```
[stable-8.2] fix(graph): [OCISDEV-1506] Enforce $filter before cross-instance user resolution
```

The marker is removed in the release notes, so the 8.2 release lists the fix with the same
words as the next master release.

### Why `-x`

`git cherry-pick -x` appends a line to the commit message:

```
(cherry picked from commit 3deec26d8933f48ac19d8ce26fd2f221e8cceebd)
```

This line is the link between the backport and the original. The `Validate backport`
check uses it to find the original PR and to verify that:

- all cherry-picked commits come from **one** PR, one backport PR per original PR
- the title is the original title, with only markers such as `[stable-8.2]` added

**Without `-x`** the backport looks like a new change. The check cannot find the original
and does not run, so a reworded or bundled backport is not caught, and the stable release
notes describe the same fix in different words or mix several fixes into one line. If you
forgot `-x`, amend the commit and add the line, or redo the cherry-pick.

**Why one PR per backport:** each PR is one line in the release notes. A backport of three
fixes in one PR becomes one line with one title, and two of the fixes disappear from the
stable release notes.

If the cherry-pick conflicts, resolve the conflicts and continue as usual, the `-x` line
is kept.

### Fixes only for a stable branch

A fix that is not needed on `master` is a normal PR against the stable branch: branch off
the stable branch, commit, open the PR with a conventional title. No cherry-pick and no
marker needed, only the `Validate PR title` check applies.

## Bots and external contributors

- **Dependabot** titles are conventional already (`build(deps): Bump ...`) and are not
  listed in the release notes. Security relevant bumps should be retitled to `security(deps): ...`
  by the reviewer.
- **External contributors** follow the same title format, a ticket is not needed. Reviewers
  may edit the title before merging, the squash commit uses the final title.

## Preview the release notes

Install [git-cliff](https://git-cliff.org/docs/installation), then:

```bash
make changelog-generate CHANGELOG_VERSION=v9.0.0
```

prints the notes for everything merged since the previous release, as the next release
`9.0.0` would get them.
