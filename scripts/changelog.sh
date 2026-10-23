#!/usr/bin/env bash
# Generate the release notes of one version from conventional commits (see
# cliff.toml), from the previous final release up to the version's tag, or up
# to HEAD when the version is not tagged yet.
#
# Usage: changelog.sh <version>            print the release notes
#        changelog.sh <version> --update   add them to CHANGELOG.md
set -euo pipefail

command -v git-cliff >/dev/null || { echo "git-cliff is not installed, see https://git-cliff.org/docs/installation" >&2; exit 1; }

version=${1:?usage: $0 <version> [--update]}
[[ $version == v* ]] || version="v$version"
file=${CHANGELOG_FILE:-CHANGELOG.md}

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

notes=$(CHANGELOG_PREVIOUS=$previous git cliff --config cliff.toml "${args[@]}" | sed '/./,$!d')

if [[ ${2:-} != --update ]]; then
  echo "$notes"
  exit 0
fi

# CHANGELOG.md: new entry at the top of the Table of Contents and the section
# above the previous newest version
v=${version#v}
if grep -qF "# Changelog for [$v]" "$file"; then
  echo "$file already contains $v" >&2
  exit 1
fi
date=$(sed -n "s/^# Changelog for \[$v\] (\(.*\))$/\1/p" <<< "$notes")
section=$(mktemp)
out=$(mktemp)
trap 'rm -f "$section" "$out"' EXIT
echo "$notes" > "$section"
awk -v toc="* [Changelog for $v](#changelog-for-${v//./}-$date)" -v section="$section" '
  /^# Table of Contents/ { print; getline; print; print toc; next }
  !done && /^# Changelog for \[/ {
    while ((getline line < section) > 0) print line
    print ""
    done = 1
  }
  { print }
' "$file" > "$out"
cp "$out" "$file"
