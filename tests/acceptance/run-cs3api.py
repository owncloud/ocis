#!/usr/bin/env python3
"""
Run CS3 API validator tests locally and in GitHub Actions CI.

Config sourced from .drone.star cs3ApiTests() — single source of truth.
Usage: python3 tests/acceptance/run-cs3api.py
"""

import os
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
CS3API_IMAGE = "owncloud/cs3api-validator:0.2.1"


def get_docker_bridge_ip() -> str:
    """Return the Docker bridge gateway IP, reachable from host and Docker containers."""
    r = subprocess.run(
        ["docker", "network", "inspect", "bridge",
         "--format", "{{range .IPAM.Config}}{{.Gateway}}{{end}}"],
        capture_output=True, text=True, check=True,
    )
    return r.stdout.strip()


def base_server_env(repo_root: Path, ocis_config_dir: str, ocis_public_url: str) -> dict:
    """OCIS server environment matching drone ocisServer(deploy_type='cs3api_validator')."""
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
        # cs3api_validator extras (drone ocisServer deploy_type="cs3api_validator")
        "GATEWAY_GRPC_ADDR": "0.0.0.0:9142",
        "OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD": "false",
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


def main() -> int:
    repo_root = Path(__file__).resolve().parents[2]
    ocis_bin = repo_root / "ocis/bin/ocis"
    ocis_config_dir = Path.home() / ".ocis/config"

    subprocess.run(["make", "-C", str(repo_root / "ocis"), "build"], check=True)

    # Docker bridge gateway IP: reachable from both the host and Docker containers.
    # cs3api-validator connects to the GRPC gateway at {bridge_ip}:9142.
    bridge_ip = get_docker_bridge_ip()
    print(f"Docker bridge IP: {bridge_ip}", flush=True)

    server_env = {**os.environ}
    server_env.update(base_server_env(repo_root, str(ocis_config_dir),
                                      f"https://{bridge_ip}:9200"))

    subprocess.run(
        [str(ocis_bin), "init", "--insecure", "true"],
        env=server_env,
        check=True,
    )
    shutil.copy(
        repo_root / "tests/config/drone/app-registry.yaml",
        ocis_config_dir / "app-registry.yaml",
    )

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

        print(f"\nRunning cs3api-validator against {bridge_ip}:9142", flush=True)
        result = subprocess.run(
            ["docker", "run", "--rm",
             "--entrypoint", "/usr/bin/cs3api-validator",
             CS3API_IMAGE,
             "/var/lib/cs3api-validator",
             f"--endpoint={bridge_ip}:9142"],
        )
        return result.returncode

    finally:
        cleanup()


if __name__ == "__main__":
    sys.exit(main())
