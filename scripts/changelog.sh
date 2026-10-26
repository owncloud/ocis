#!/usr/bin/env bash
# Print the release notes of one version from conventional commits (see
# cliff.toml), from the previous final release up to the version's tag, or up
# to HEAD when the version is not tagged yet.
#
# Usage: changelog.sh <version>
set -euo pipefail

command -v git-cliff >/dev/null || { echo "git-cliff is not installed, see https://git-cliff.org/docs/installation" >&2; exit 1; }

version=${1:?usage: $0 <version>}
[[ $version == v* ]] || version="v$version"

# previous final release by version, so 8.2.1 follows 8.2.0 and 9.0.0 follows
# the newest 8.x; rc, dev and web package tags are ignored
previous=$( { git tag -l 'v[0-9]*' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$'; echo "$version"; } \
  | sort -uV | grep -B1 -Fx "$version" | grep -vFx "$version" || true)

# up to the tag, or up to HEAD named as the version when it is not tagged yet
if git rev-parse -q --verify "refs/tags/$version" >/dev/null; then
  args=("${previous:+$previous..}$version")
else
  args=("${previous:+$previous..}HEAD" --tag "$version")
fi

CHANGELOG_PREVIOUS=$previous git cliff --config cliff.toml "${args[@]}" | sed '/./,$!d'
