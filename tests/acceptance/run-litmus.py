#!/usr/bin/env python3
"""
Run litmus WebDAV compliance tests locally and in GitHub Actions CI.

Config sourced from .drone.star litmus() / setupForLitmus() — single source of truth.
Usage: python3 tests/acceptance/run-litmus.py
"""

import json
import os
import re
import shutil
import signal
import subprocess
import sys
import time
from pathlib import Path

# ---------------------------------------------------------------------------
# Constants (mirroring .drone.star)
# ---------------------------------------------------------------------------

OCIS_URL = "http://127.0.0.1:9200"
LITMUS_IMAGE = "owncloudci/litmus:latest"
LITMUS_TESTS = "basic copymove props http"
SHARE_ENDPOINT = "ocs/v2.php/apps/files_sharing/api/v1/shares"


def get_docker_bridge_ip() -> str:
    """Return the Docker bridge gateway IP, reachable from host and Docker containers."""
    r = subprocess.run(
        ["docker", "network", "inspect", "bridge",
         "--format", "{{range .IPAM.Config}}{{.Gateway}}{{end}}"],
        capture_output=True, text=True, check=True,
    )
    return r.stdout.strip()


def base_server_env(repo_root: Path, ocis_config_dir: str, ocis_public_url: str) -> dict:
    """OCIS server environment matching drone ocisServer() for litmus."""
    return {
        "OCIS_URL": ocis_public_url,
        "OCIS_CONFIG_DIR": ocis_config_dir,
        "STORAGE_USERS_DRIVER": "ocis",
        "PROXY_ENABLE_BASIC_AUTH": "true",
        "PROXY_TLS": "false",
        "OCIS_EXCLUDE_RUN_SERVICES": "idp",
        "OCIS_LOG_LEVEL": "error",
        "IDM_CREATE_DEMO_USERS": "true",
        "IDM_ADMIN_PASSWORD": "admin",
        "FRONTEND_SEARCH_MIN_LENGTH": "2",
        "OCIS_ASYNC_UPLOADS": "true",
        "OCIS_EVENTS_ENABLE_TLS": "false",
        "NATS_NATS_HOST": "0.0.0.0",
        "NATS_NATS_PORT": "9233",
        "OCIS_JWT_SECRET": "some-ocis-jwt-secret",
        "EVENTHISTORY_STORE": "memory",
        "WEB_UI_CONFIG_FILE": str(repo_root / "tests/config/drone/ocis-config.json"),
    }


def wait_for(condition_fn, timeout: int, label: str) -> None:
    deadline = time.time() + timeout
    while not condition_fn():
        if time.time() > deadline:
            print(f"Timeout waiting for {label}", file=sys.stderr)
            sys.exit(1)
        time.sleep(1)


def ocis_healthy(ocis_url: str) -> bool:
    r = subprocess.run(
        ["curl", "-s", "-uadmin:admin",
         f"{ocis_url}/graph/v1.0/users/admin",
         "-w", "%{http_code}", "-o", "/dev/null"],
        capture_output=True, text=True,
    )
    return r.stdout.strip() == "200"


def setup_for_litmus(ocis_url: str) -> tuple:
    """
    Translate tests/config/drone/setup-for-litmus.sh to Python.
    Returns (space_id, public_token).
    """
    # get personal space ID
    r = subprocess.run(
        ["curl", "-s", "-uadmin:admin", f"{ocis_url}/graph/v1.0/me/drives"],
        capture_output=True, text=True, check=True,
    )
    drives = json.loads(r.stdout)
    space_id = ""
    for drive in drives.get("value", []):
        if drive.get("driveType") == "personal":
            web_dav_url = drive.get("root", {}).get("webDavUrl", "")
            # last non-empty path segment (same as cut -d"/" -f6 in bash)
            space_id = [p for p in web_dav_url.split("/") if p][-1]
            break
    if not space_id:
        print("ERROR: could not determine personal space ID", file=sys.stderr)
        sys.exit(1)
    print(f"SPACE_ID={space_id}")

    # create test folder as einstein
    subprocess.run(
        ["curl", "-s", "-ueinstein:relativity", "-X", "MKCOL",
         f"{ocis_url}/remote.php/webdav/new_folder"],
        capture_output=True, check=True,
    )

    # create share from einstein to admin
    r = subprocess.run(
        ["curl", "-s", "-ueinstein:relativity",
         f"{ocis_url}/{SHARE_ENDPOINT}",
         "-d", "path=/new_folder&shareType=0&permissions=15&name=new_folder&shareWith=admin"],
        capture_output=True, text=True, check=True,
    )
    share_id_match = re.search(r"<id>(.+?)</id>", r.stdout)
    if share_id_match:
        share_id = share_id_match.group(1)
        # accept the share as admin
        subprocess.run(
            ["curl", "-X", "POST", "-s", "-uadmin:admin",
             f"{ocis_url}/{SHARE_ENDPOINT}/pending/{share_id}"],
            capture_output=True, check=True,
        )

    # create public share as einstein
    r = subprocess.run(
        ["curl", "-s", "-ueinstein:relativity",
         f"{ocis_url}/{SHARE_ENDPOINT}",
         "-d", "path=/new_folder&shareType=3&permissions=15&name=new_folder"],
        capture_output=True, text=True, check=True,
    )
    public_token = ""
    token_match = re.search(r"<token>(.+?)</token>", r.stdout)
    if token_match:
        public_token = token_match.group(1)
    print(f"PUBLIC_TOKEN={public_token}")

    return space_id, public_token


def run_litmus(name: str, endpoint: str) -> int:
    print(f"\nTesting endpoint [{name}]: {endpoint}", flush=True)
    result = subprocess.run(
        ["docker", "run", "--rm",
         "-e", f"LITMUS_URL={endpoint}",
         "-e", "LITMUS_USERNAME=admin",
         "-e", "LITMUS_PASSWORD=admin",
         "-e", f"TESTS={LITMUS_TESTS}",
         LITMUS_IMAGE,
         "/usr/local/bin/litmus-wrapper"],
    )
    return result.returncode


def main() -> int:
    repo_root = Path(__file__).resolve().parents[2]
    ocis_bin = repo_root / "ocis/bin/ocis"
    ocis_config_dir = Path.home() / ".ocis/config"

    # build (matching drone: restores binary from cache, then runs ocis server directly)
    subprocess.run(["make", "-C", str(repo_root / "ocis"), "build"], check=True)

    # Docker bridge gateway IP: reachable from both the host (via docker0 interface)
    # and Docker containers (via bridge network default gateway). Use this as OCIS_URL
    # so that any redirects OCIS generates stay on a hostname the litmus container
    # can follow — matching how drone uses "ocis-server:9200" consistently.
    bridge_ip = get_docker_bridge_ip()
    litmus_base = f"http://{bridge_ip}:9200"
    print(f"Docker bridge IP: {bridge_ip}", flush=True)

    # assemble server env first — same env vars drone sets on the container before
    # running `ocis init`, so IDM_ADMIN_PASSWORD=admin is present during init and
    # the config is written with the correct password (not a random one)
    server_env = {**os.environ}
    server_env.update(base_server_env(repo_root, str(ocis_config_dir), litmus_base))

    # init ocis with full server env (mirrors drone: env is set before ocis init runs)
    subprocess.run(
        [str(ocis_bin), "init", "--insecure", "true"],
        env=server_env,
        check=True,
    )
    shutil.copy(
        repo_root / "tests/config/drone/app-registry.yaml",
        ocis_config_dir / "app-registry.yaml",
    )

    # start ocis server directly (matching drone: no ociswrapper for litmus)
    print("Starting ocis...", flush=True)
    ocis_proc = subprocess.Popen(
        [str(ocis_bin), "server"],
        env=server_env,
    )

    def cleanup(*_):
        try:
            ocis_proc.terminate()
        except Exception:
            pass

    signal.signal(signal.SIGTERM, cleanup)
    signal.signal(signal.SIGINT, cleanup)

    try:
        wait_for(lambda: ocis_healthy(OCIS_URL), 300, "ocis")
        print("ocis ready.", flush=True)

        space_id, _ = setup_for_litmus(OCIS_URL)

        # Diagnostic: test TCP + HTTP connectivity to OCIS from Docker containers
        # BusyBox wget --spider makes a GET and reports HTTP status without downloading body
        for label, extra_args, test_url in [
            ("docker-host-net", ["--network", "host"], "http://localhost:9200/graph/v1.0/users/admin"),
            ("docker-bridge",   [],                    f"http://{bridge_ip}:9200/graph/v1.0/users/admin"),
        ]:
            wget_cmd = f"wget -S --spider '{test_url}' 2>&1; echo exit:$?"
            r = subprocess.run(
                ["docker", "run", "--rm"] + extra_args +
                ["--entrypoint", "sh", LITMUS_IMAGE, "-c", wget_cmd],
                capture_output=True, text=True, timeout=30,
            )
            print(f"\n--- {label} rc={r.returncode} ---", flush=True)
            print((r.stdout + r.stderr)[:400], flush=True)

        endpoints = [
            ("old-endpoint",    f"{litmus_base}/remote.php/webdav"),
            ("new-endpoint",    f"{litmus_base}/remote.php/dav/files/admin"),
            ("new-shared",      f"{litmus_base}/remote.php/dav/files/admin/Shares/new_folder/"),
            ("old-shared",      f"{litmus_base}/remote.php/webdav/Shares/new_folder/"),
            ("spaces-endpoint", f"{litmus_base}/remote.php/dav/spaces/{space_id}"),
        ]

        failed = []
        for name, endpoint in endpoints:
            rc = run_litmus(name, endpoint)
            if rc != 0:
                failed.append(name)

        if failed:
            print(f"\nFailed endpoints: {', '.join(failed)}", file=sys.stderr)
            return 1
        print("\nAll litmus tests passed.")
        return 0

    finally:
        cleanup()


if __name__ == "__main__":
    sys.exit(main())
