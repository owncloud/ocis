#!/usr/bin/env python3
import argparse
import hashlib
import json
import os
import struct
import subprocess
import sys
import urllib.error
import urllib.request
from pathlib import Path

BINARY_PLATFORMS = ["darwin-amd64", "darwin-arm64", "linux-386", "linux-amd64", "linux-arm", "linux-arm64"]
MACHO_MAGIC      = b"\xcf\xfa\xed\xfe"
ELF_MAGIC        = b"\x7fELF"
BINARY_REF = {
    "linux-amd64":  (ELF_MAGIC,   64, 0x3e),
    "linux-arm64":  (ELF_MAGIC,   64, 0xb7),
    "linux-arm":    (ELF_MAGIC,   32, 0x28),
    "linux-386":    (ELF_MAGIC,   32, 0x03),
    "darwin-amd64": (MACHO_MAGIC, 64, 0x01000007),
    "darwin-arm64": (MACHO_MAGIC, 64, 0x0100000c),
}
EULA      = "End-User-License-Agreement-for-ownCloud-Infinite-Scale.pdf"
LICENSES  = "third-party-licenses.tar.gz"
PROD_TAGS = ("5.0", "7", "8")

passed = failed = 0

def ok(msg):
    global passed; passed += 1; print(f"[PASS] {msg}")

def fail(msg):
    global failed; failed += 1; print(f"[FAIL] {msg}", file=sys.stderr)

def expected_files(v):
    return [f for p in BINARY_PLATFORMS for f in (f"ocis-{v}-{p}", f"ocis-{v}-{p}.sha256")] + [EULA, LICENSES]

def is_production(v):
    return any(v.startswith(t) for t in PROD_TAGS) and "-" not in v

def run(cmd):
    return subprocess.run(cmd, capture_output=True, text=True)

def get_github_token():
    return os.environ.get("GH_TOKEN") or os.environ.get("GITHUB_TOKEN")

def gh_api(path, token):
    req = urllib.request.Request(
        f"https://api.github.com{path}",
        headers={"Authorization": f"Bearer {token}", "Accept": "application/vnd.github+json"},
    )
    with urllib.request.urlopen(req) as r:
        return json.load(r)

def check_file_magic(path, expected_magic, min_size, fmt_name):
    if not path.exists():
        return
    data = path.read_bytes()
    if not data.startswith(expected_magic):
        fail(f"{path.name}: not {fmt_name} ({data[:len(expected_magic)].hex()})")
    elif len(data) < min_size:
        fail(f"{path.name}: too small ({len(data):,} bytes)")
    else:
        ok(f"{path.name}: {fmt_name} {len(data):,} bytes")


def check_local(directory, version):
    expected = expected_files(version)
    present  = {f.name for f in directory.iterdir()}
    missing  = [f for f in expected if f not in present]
    extra    = present - set(expected)
    if missing: fail(f"file set: missing {missing}")
    else:       ok(f"file set: {len(expected)} files present")
    if extra:   fail(f"file set: unexpected {extra}")

    for platform in BINARY_PLATFORMS:
        bin_path  = directory / f"ocis-{version}-{platform}"
        sha_path  = directory / f"ocis-{version}-{platform}.sha256"
        bin_bytes = bin_path.read_bytes() if bin_path.exists() else None

        if bin_bytes is not None:
            magic, bits, machine = BINARY_REF[platform]
            if len(bin_bytes) == 0:
                fail(f"{bin_path.name}: empty file")
            else:
                h = bin_bytes[:20]
                if not h.startswith(magic):
                    fail(f"{bin_path.name}: wrong magic {h[:4].hex()}")
                elif magic == ELF_MAGIC:
                    cls, mach = h[4], struct.unpack_from("<H", h, 18)[0]
                    if (64 if cls == 2 else 32) != bits or mach != machine:
                        fail(f"{bin_path.name}: wrong ELF class={cls} machine=0x{mach:x}")
                    else:
                        ok(f"{bin_path.name}: ELF {bits}-bit 0x{mach:x}")
                else:
                    cputype = struct.unpack_from("<I", h, 4)[0]
                    if cputype != machine:
                        fail(f"{bin_path.name}: wrong Mach-O cputype 0x{cputype:08x}")
                    else:
                        ok(f"{bin_path.name}: Mach-O 0x{cputype:08x}")

        if sha_path.exists():
            if bin_bytes is None:
                fail(f"{sha_path.name}: binary missing")
                continue
            parts = sha_path.read_text().strip().split("  ", 1)
            if len(parts) != 2:
                fail(f"{sha_path.name}: bad format")
                continue
            rec_hash, rec_name = parts
            if rec_name != bin_path.name:
                fail(f"{sha_path.name}: filename mismatch '{rec_name}'")
            actual = hashlib.sha256(bin_bytes).hexdigest()
            if actual != rec_hash:
                fail(f"{sha_path.name}: hash mismatch\n  want: {rec_hash}\n  got:  {actual}")
            else:
                ok(f"{sha_path.name}: hash ok")

    check_file_magic(directory / LICENSES, b"\x1f\x8b", 100_000, "gzip")
    check_file_magic(directory / EULA,     b"%PDF",     10_000,  "PDF")


def check_github_release(version):
    token = get_github_token()
    if not token:
        fail("GH_TOKEN / GITHUB_TOKEN not set"); return
    try:
        r = gh_api(f"/repos/owncloud/ocis/releases/tags/v{version}", token)
    except urllib.error.HTTPError as e:
        fail(f"release v{version} not found: {e}"); return

    checks = [
        (r.get("tag_name") == f"v{version}", f"tag_name: {r.get('tag_name')}"),
        (r.get("name") == f"v{version}",       f"name: {r.get('name')}"),
        (not r.get("draft"),                  "draft: false"),
        (r.get("prerelease") == ("-" in version), f"prerelease: {r.get('prerelease')}"),
    ]
    for result, label in checks:
        ok(label) if result else fail(label)

    published = {a["name"] for a in r.get("assets", [])}
    expected  = set(expected_files(version))
    missing, extra = expected - published, published - expected
    if missing: fail(f"assets missing: {sorted(missing)}")
    else:       ok(f"assets: all {len(expected)} present")
    if extra:   fail(f"assets unexpected: {sorted(extra)}")


def check_docker(version):
    # Returns dict of "os/arch" -> digest for a multi-arch manifest, or (None, err).
    def inspect_manifest(ref):
        r = run(["docker", "buildx", "imagetools", "inspect", "--format",
                 "{{range .Manifest.Manifests}}{{.Platform.OS}}/{{.Platform.Architecture}}={{.Digest}} {{end}}", ref])
        if r.returncode != 0:
            return None, r.stderr.strip()
        return dict(e.split("=", 1) for e in r.stdout.split() if "=" in e), None

    def check_arches(ref, manifests):
        missing = [a for a in ("linux/amd64", "linux/arm64") if a not in manifests]
        if missing: fail(f"{ref}: missing {missing}")
        else:       ok(f"{ref}: amd64+arm64 present")

    rolling = f"owncloud/ocis-rolling:{version}"
    m, err = inspect_manifest(rolling)
    if m is None: fail(f"{rolling}: {err}")
    else:         check_arches(rolling, m)

    if not is_production(version):
        return

    prod = f"owncloud/ocis:{version}"
    versioned, err = inspect_manifest(prod)
    if versioned is None:
        fail(f"{prod}: {err}")
    else:
        check_arches(prod, versioned)
    versioned_digests = set(versioned.values()) if versioned else set()

    parts = version.split(".")
    major, major_minor = parts[0], ".".join(parts[:2])
    # If a newer minor/major is already live (e.g. auditing an 8.0.x backport
    # while 8.1 is out) floating tags must NOT be downgraded — that's a regression.
    for floating in (f"owncloud/ocis:{major_minor}", f"owncloud/ocis:{major}", "owncloud/ocis:latest"):
        m, err = inspect_manifest(floating)
        if m is None:
            fail(f"{floating}: not found or inspect failed — tag was not pushed")
            continue
        if not versioned_digests:
            fail(f"{floating}: cannot compare — versioned manifest inspect failed")
            continue
        if set(m.values()) == versioned_digests:
            ok(f"{floating}: matches {version}")
            continue
        fail(f"{floating}: points to a different manifest than {version} — floating tag not updated or downgrade regression")


def resolve_run_id(branch, token):
    r = gh_api(f"/repos/owncloud/ocis/actions/workflows/release.yml/runs?branch={branch}&per_page=1", token)
    runs = r.get("workflow_runs", [])
    if not runs:
        sys.exit(f"no runs found for branch '{branch}'")
    run = runs[0]
    print(f"run {run['id']}  {run['status']}  {run['conclusion'] or 'in_progress'}  {run['html_url']}")
    return str(run["id"])


def check_run_artifacts(run_id):
    token = get_github_token()
    if not token:
        fail("GH_TOKEN / GITHUB_TOKEN not set"); return
    try:
        r = gh_api(f"/repos/owncloud/ocis/actions/runs/{run_id}/artifacts", token)
    except urllib.error.HTTPError as e:
        fail(f"run {run_id}: {e}"); return

    by_name = {a["name"]: a for a in r.get("artifacts", [])}
    for name in ("binaries-linux", "binaries-darwin", "third-party-licenses"):
        a = by_name.get(name)
        if not a:                   fail(f"{name}: missing");  continue
        if a.get("expired"):        fail(f"{name}: expired");  continue
        if a["size_in_bytes"] == 0: fail(f"{name}: empty");    continue
        ok(f"{name}: {a['size_in_bytes']:,} bytes")


def check_git(version):
    tag = f"v{version}"
    r = run(["git", "tag", "-v", tag])
    if r.returncode != 0: print(f"[WARN] {tag}: not a signed tag")
    else:                 ok(f"{tag}: signed tag ok")

    r = run(["git", "cat-file", "-p", tag])
    if r.returncode != 0:
        fail(f"{tag}: cat-file failed")
    else:
        obj = next((l for l in r.stdout.splitlines() if l.startswith("object")), None)
        ok(f"{tag}: {obj}") if obj else fail(f"{tag}: no object line in cat-file")


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--version",        required=False)
    p.add_argument("--run",            metavar="RUN_ID")
    p.add_argument("--branch",         metavar="BRANCH")
    p.add_argument("--dir",            type=Path)
    p.add_argument("--github-release", action="store_true")
    p.add_argument("--docker",         action="store_true")
    p.add_argument("--git",            action="store_true")
    args = p.parse_args()

    if not args.run and not args.branch and not args.version:
        p.error("--version is required unless --run or --branch is used")

    if args.branch:
        token = get_github_token()
        if not token: sys.exit("GH_TOKEN / GITHUB_TOKEN not set")
        args.run = resolve_run_id(args.branch, token)

    needs_version = args.dir or args.github_release or args.docker or args.git
    if needs_version and not args.version:
        p.error("--version is required with --dir / --github-release / --docker / --git")

    if args.run:             check_run_artifacts(args.run)
    if args.dir:             check_local(args.dir, args.version)
    if args.github_release:  check_github_release(args.version)
    if args.docker:          check_docker(args.version)
    if args.git:             check_git(args.version)

    print(f"\n{passed + failed} checks: {passed} passed, {failed} failed")
    sys.exit(1 if failed else 0)


if __name__ == "__main__":
    main()
