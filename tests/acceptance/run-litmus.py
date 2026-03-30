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

# HTTPS — matching drone: ocis init generates a self-signed cert; proxy uses TLS by default.
# Host-side curl calls use -k (insecure) to skip cert verification.
OCIS_URL = "https://127.0.0.1:9200"
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
        # No PROXY_TLS override — drone lets ocis use its default TLS (self-signed cert from init)
        # IDP excluded: its static assets are absent when running as a host process
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
        ["curl", "-sk", "-uadmin:admin",
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
        ["curl", "-sk", "-uadmin:admin", f"{ocis_url}/graph/v1.0/me/drives"],
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
        ["curl", "-sk", "-ueinstein:relativity", "-X", "MKCOL",
         f"{ocis_url}/remote.php/webdav/new_folder"],
        capture_output=True, check=True,
    )

    # create share from einstein to admin
    r = subprocess.run(
        ["curl", "-sk", "-ueinstein:relativity",
         f"{ocis_url}/{SHARE_ENDPOINT}",
         "-d", "path=/new_folder&shareType=0&permissions=15&name=new_folder&shareWith=admin"],
        capture_output=True, text=True, check=True,
    )
    share_id_match = re.search(r"<id>(.+?)</id>", r.stdout)
    if share_id_match:
        share_id = share_id_match.group(1)
        # accept the share as admin
        subprocess.run(
            ["curl", "-X", "POST", "-sk", "-uadmin:admin",
             f"{ocis_url}/{SHARE_ENDPOINT}/pending/{share_id}"],
            capture_output=True, check=True,
        )

    # create public share as einstein
    r = subprocess.run(
        ["curl", "-sk", "-ueinstein:relativity",
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


def run_litmus(name: str, endpoint: str, capture_debug: bool = False) -> int:
    print(f"\nTesting endpoint [{name}]: {endpoint}", flush=True)
    container_name = f"litmus-{name}" if capture_debug else None
    cmd = ["docker", "run"]
    if container_name:
        cmd += ["--name", container_name]
    else:
        cmd += ["--rm"]
    cmd += [
        "-e", f"LITMUS_URL={endpoint}",
        "-e", "LITMUS_USERNAME=admin",
        "-e", "LITMUS_PASSWORD=admin",
        "-e", f"TESTS={LITMUS_TESTS}",
        LITMUS_IMAGE,
        # No extra CMD — ENTRYPOINT is already litmus-wrapper; passing it again
        # would make the wrapper use the path as LITMUS_URL, overriding the env var.
    ]
    result = subprocess.run(cmd)
    if capture_debug and container_name:
        # try to copy debug.log from the container
        for log_path in ["/debug.log", "/tmp/debug.log", "/home/debug.log"]:
            r = subprocess.run(
                ["docker", "cp", f"{container_name}:{log_path}", f"/tmp/litmus-debug-{name}.log"],
                capture_output=True,
            )
            if r.returncode == 0:
                try:
                    with open(f"/tmp/litmus-debug-{name}.log") as f:
                        content = f.read()
                    print(f"\n--- litmus debug.log ({log_path}) ---", flush=True)
                    print(content[:3000], flush=True)
                except Exception:
                    pass
                break
        subprocess.run(["docker", "rm", "-f", container_name], capture_output=True)
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
    # HTTPS: owncloudci/litmus accepts insecure (self-signed) certs, just like drone does.
    bridge_ip = get_docker_bridge_ip()
    litmus_base = f"https://{bridge_ip}:9200"
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
    # Diagnostic: print generated config (drone does this too)
    ocis_yaml = ocis_config_dir / "ocis.yaml"
    if ocis_yaml.exists():
        print("\n--- ocis.yaml ---", flush=True)
        print(ocis_yaml.read_text()[:2000], flush=True)

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

        # Diagnostic: full HTTP trace for WebDAV OPTIONS and PROPFIND from the host
        for label, curl_args in [
            ("OPTIONS no-auth",  ["-X", "OPTIONS"]),
            ("OPTIONS auth",     ["-u", "admin:admin", "-X", "OPTIONS"]),
            ("PROPFIND auth",    ["-u", "admin:admin", "-X", "PROPFIND", "-H", "Depth: 0"]),
        ]:
            r = subprocess.run(
                ["curl", "-vsk"] + curl_args + [f"{OCIS_URL}/remote.php/webdav"],
                capture_output=True, text=True,
            )
            print(f"\n--- {label} ---", flush=True)
            print((r.stdout + r.stderr)[:1500], flush=True)

        # Diagnostic: confirm WebDAV endpoint reachable from Docker bridge network (HTTPS, insecure)
        webdav_url = f"https://{bridge_ip}:9200/remote.php/webdav"
        wget_cmd = f"wget -S --no-check-certificate --spider '{webdav_url}' 2>&1; echo exit:$?"
        r = subprocess.run(
            ["docker", "run", "--rm", "--entrypoint", "sh", LITMUS_IMAGE, "-c", wget_cmd],
            capture_output=True, text=True, timeout=30,
        )
        print(f"\n--- docker-bridge WebDAV (https) rc={r.returncode} ---", flush=True)
        print((r.stdout + r.stderr)[:400], flush=True)

        endpoints = [
            ("old-endpoint",    f"{litmus_base}/remote.php/webdav"),
            ("new-endpoint",    f"{litmus_base}/remote.php/dav/files/admin"),
            ("new-shared",      f"{litmus_base}/remote.php/dav/files/admin/Shares/new_folder/"),
            ("old-shared",      f"{litmus_base}/remote.php/webdav/Shares/new_folder/"),
            ("spaces-endpoint", f"{litmus_base}/remote.php/dav/spaces/{space_id}"),
        ]

        failed = []
        for i, (name, endpoint) in enumerate(endpoints):
            rc = run_litmus(name, endpoint, capture_debug=(i == 0))
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
