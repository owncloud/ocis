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
    return any(v.startswith(t) for t in PROD_TAGS)

def run(cmd):
    return subprocess.run(cmd, capture_output=True, text=True)

def gh_api(path, token):
    req = urllib.request.Request(
        f"https://api.github.com{path}",
        headers={"Authorization": f"Bearer {token}", "Accept": "application/vnd.github+json"},
    )
    with urllib.request.urlopen(req) as r:
        return json.load(r)


def check_local(directory, version):
    expected = expected_files(version)
    present  = {f.name for f in directory.iterdir()}
    missing  = [f for f in expected if f not in present]
    extra    = present - set(expected)
    if missing: fail(f"file set: missing {missing}")
    else:       ok(f"file set: {len(expected)} files present")
    if extra:   fail(f"file set: unexpected {extra}")

    for platform in BINARY_PLATFORMS:
        bin_path = directory / f"ocis-{version}-{platform}"
        sha_path = directory / f"ocis-{version}-{platform}.sha256"

        if bin_path.exists():
            magic, bits, machine = BINARY_REF[platform]
            if bin_path.stat().st_size == 0:
                fail(f"{bin_path.name}: empty file")
            else:
                h = bin_path.read_bytes()[:20]
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
            if not bin_path.exists():
                fail(f"{sha_path.name}: binary missing")
                continue
            size = sha_path.stat().st_size
            if not (87 <= size <= 90):
                fail(f"{sha_path.name}: size {size} (want 87-90)")
            parts = sha_path.read_text().strip().split("  ", 1)
            if len(parts) != 2:
                fail(f"{sha_path.name}: bad format")
                continue
            rec_hash, rec_name = parts
            if rec_name != bin_path.name:
                fail(f"{sha_path.name}: filename mismatch '{rec_name}'")
            actual = hashlib.sha256(bin_path.read_bytes()).hexdigest()
            if actual != rec_hash:
                fail(f"{sha_path.name}: hash mismatch\n  want: {rec_hash}\n  got:  {actual}")
            else:
                ok(f"{sha_path.name}: hash ok")

    lic = directory / LICENSES
    if lic.exists():
        magic, size = lic.read_bytes()[:2], lic.stat().st_size
        if magic != b"\x1f\x8b": fail(f"{LICENSES}: not gzip ({magic.hex()})")
        elif size < 100_000:     fail(f"{LICENSES}: too small ({size:,} bytes)")
        else:                    ok(f"{LICENSES}: gzip {size:,} bytes")

    eula = directory / EULA
    if eula.exists():
        magic, size = eula.read_bytes()[:4], eula.stat().st_size
        if magic != b"%PDF":  fail(f"{EULA}: not PDF ({magic})")
        elif size < 10_000:   fail(f"{EULA}: too small ({size:,} bytes)")
        else:                 ok(f"{EULA}: PDF {size:,} bytes")


def check_github_release(version):
    token = os.environ.get("GH_TOKEN") or os.environ.get("GITHUB_TOKEN")
    if not token:
        fail("GH_TOKEN / GITHUB_TOKEN not set"); return
    try:
        r = gh_api(f"/repos/owncloud/ocis/releases/tags/v{version}", token)
    except urllib.error.HTTPError as e:
        fail(f"release v{version} not found: {e}"); return

    checks = [
        (r.get("tag_name") == f"v{version}", f"tag_name: {r.get('tag_name')}"),
        (r.get("name") == version,            f"name: {r.get('name')}"),
        (not r.get("draft"),                  "draft: false"),
        (r.get("prerelease") == ("-" in version), f"prerelease: {r.get('prerelease')}"),
    ]
    for passed_check, label in checks:
        ok(label) if passed_check else fail(label)

    published = {a["name"] for a in r.get("assets", [])}
    expected  = set(expected_files(version))
    missing, extra = expected - published, published - expected
    if missing: fail(f"assets missing: {sorted(missing)}")
    else:       ok(f"assets: all {len(expected)} present")
    if extra:   fail(f"assets unexpected: {sorted(extra)}")


def check_docker(version):
    refs = [f"owncloud/ocis-rolling:{version}"]
    if is_production(version):
        refs.append(f"owncloud/ocis:{version}")
    for ref in refs:
        r = run(["docker", "buildx", "imagetools", "inspect", ref])
        if r.returncode != 0:
            fail(f"{ref}: {r.stderr.strip()}"); continue
        missing = [a for a in ("linux/amd64", "linux/arm64") if a not in r.stdout]
        if missing: fail(f"{ref}: missing {missing}")
        else:       ok(f"{ref}: amd64+arm64 present")


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
    p.add_argument("--version",        required=True)
    p.add_argument("--dir",            type=Path)
    p.add_argument("--github-release", action="store_true")
    p.add_argument("--docker",         action="store_true")
    p.add_argument("--git",            action="store_true")
    args = p.parse_args()

    if args.dir:            check_local(args.dir, args.version)
    if args.github_release: check_github_release(args.version)
    if args.docker:         check_docker(args.version)
    if args.git:            check_git(args.version)

    print(f"\n{passed + failed} checks: {passed} passed, {failed} failed")
    sys.exit(1 if failed else 0)


if __name__ == "__main__":
    main()
