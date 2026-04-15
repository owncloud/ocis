#!/usr/bin/env python3
"""
Run WOPI validator tests locally and in GitHub Actions CI.

Config sourced from .drone.star wopiValidatorTests() — single source of truth.
TODO: [DRONE-REMOVAL] decouple from .drone.star constants when Drone CI is removed.
Usage: python3 tests/acceptance/run-wopi.py --type builtin
       python3 tests/acceptance/run-wopi.py --type cs3
"""

import argparse
import json
import os
import re
import shutil
import signal
import socket
import subprocess
import sys
import time
import urllib.parse
from pathlib import Path

# ---------------------------------------------------------------------------
# Constants (mirroring .drone.star)  # TODO: [DRONE-REMOVAL]
# ---------------------------------------------------------------------------

# HTTPS — matching drone; host-side curl calls use -k.
OCIS_URL = "https://127.0.0.1:9200"
VALIDATOR_IMAGE = "owncloudci/wopi-validator"
CS3_WOPI_IMAGE = "cs3org/wopiserver:v10.4.0"
FAKEOFFICE_IMAGE = "owncloudci/alpine:latest"

# Testgroups shared between both variants (drone: testgroups list)
SHARED_TESTGROUPS = [
    "BaseWopiViewing",
    "CheckFileInfoSchema",
    "EditFlows",
    "Locks",
    "AccessTokens",
    "GetLock",
    "ExtendedLockLength",
    "FileVersion",
    "Features",
]

# Testgroups only run for builtin (drone: builtinOnlyTestGroups, with -s flag)
BUILTIN_ONLY_TESTGROUPS = [
    "PutRelativeFile",
    "RenameFileIfCreateChildFileIsNotSupported",
]


def get_docker_bridge_ip() -> str:
    r = subprocess.run(
        ["docker", "network", "inspect", "bridge",
         "--format", "{{range .IPAM.Config}}{{.Gateway}}{{end}}"],
        capture_output=True, text=True, check=True,
    )
    return r.stdout.strip()


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


def tcp_reachable(host: str, port: int) -> bool:
    try:
        with socket.create_connection((host, port), timeout=1):
            return True
    except Exception:
        return False


def wopi_discovery_ready(url: str) -> bool:
    r = subprocess.run(
        ["curl", "-sk", "-o", "/dev/null", "-w", "%{http_code}", url],
        capture_output=True, text=True,
    )
    return r.stdout.strip() == "200"


def base_server_env(repo_root: Path, ocis_config_dir: str, ocis_public_url: str,
                    bridge_ip: str, wopi_type: str) -> dict:
    """
    OCIS server environment matching drone ocisServer(deploy_type='wopi_validator').
    builtin: also excludes app-provider (collaboration service takes that role).
    """
    exclude = "idp,app-provider" if wopi_type == "builtin" else "idp"
    return {
        "OCIS_URL": ocis_public_url,
        "OCIS_CONFIG_DIR": ocis_config_dir,
        "STORAGE_USERS_DRIVER": "ocis",
        "PROXY_ENABLE_BASIC_AUTH": "true",
        # No PROXY_TLS override — drone uses default TLS (self-signed cert from init)
        # IDP excluded: static assets absent when running as host process
        "OCIS_EXCLUDE_RUN_SERVICES": exclude,
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
        # wopi_validator extras (drone ocisServer deploy_type="wopi_validator")
        "GATEWAY_GRPC_ADDR": "0.0.0.0:9142",
        "APP_PROVIDER_EXTERNAL_ADDR": "com.owncloud.api.app-provider",
        "APP_PROVIDER_DRIVER": "wopi",
        "APP_PROVIDER_WOPI_APP_NAME": "FakeOffice",
        "APP_PROVIDER_WOPI_APP_URL": f"http://{bridge_ip}:8080",
        "APP_PROVIDER_WOPI_INSECURE": "true",
        "APP_PROVIDER_WOPI_WOPI_SERVER_EXTERNAL_URL": f"http://{bridge_ip}:9300",
        "APP_PROVIDER_WOPI_FOLDER_URL_BASE_URL": ocis_public_url,
    }


def collab_service_env(bridge_ip: str, ocis_config_dir: str) -> dict:
    """
    Environment for 'ocis collaboration server' (builtin wopi-fakeoffice).
    Mirrors drone wopiCollaborationService("fakeoffice").
    """
    return {
        "OCIS_URL": f"https://{bridge_ip}:9200",
        "OCIS_CONFIG_DIR": ocis_config_dir,
        "MICRO_REGISTRY": "nats-js-kv",
        "MICRO_REGISTRY_ADDRESS": "127.0.0.1:9233",
        "COLLABORATION_LOG_LEVEL": "debug",
        "COLLABORATION_GRPC_ADDR": "0.0.0.0:9301",
        "COLLABORATION_HTTP_ADDR": "0.0.0.0:9300",
        "COLLABORATION_DEBUG_ADDR": "0.0.0.0:9304",
        "COLLABORATION_APP_PROOF_DISABLE": "true",
        "COLLABORATION_APP_INSECURE": "true",
        "COLLABORATION_CS3API_DATAGATEWAY_INSECURE": "true",
        "OCIS_JWT_SECRET": "some-ocis-jwt-secret",
        "COLLABORATION_WOPI_SECRET": "some-wopi-secret",
        "COLLABORATION_APP_NAME": "FakeOffice",
        "COLLABORATION_APP_PRODUCT": "Microsoft",
        "COLLABORATION_APP_ADDR": f"http://{bridge_ip}:8080",
        # COLLABORATION_WOPI_SRC is what OCIS tells clients to use — must be reachable
        # from Docker validator containers (collaboration service runs as host process)
        "COLLABORATION_WOPI_SRC": f"http://{bridge_ip}:9300",
    }


def prepare_test_file(bridge_ip: str) -> tuple:
    """
    Upload test.wopitest via WebDAV, open the WOPI app, extract credentials.
    Mirrors the prepare-test-file step from drone.star.
    Returns (access_token, access_token_ttl, wopi_src).
    """
    headers_file = "/tmp/wopi-headers.txt"

    # PUT empty test file (--retry-connrefused/--retry-all-errors matching drone)
    subprocess.run(
        ["curl", "-sk", "-u", "admin:admin", "-X", "PUT",
         "--fail", "--retry-connrefused", "--retry", "7", "--retry-all-errors",
         f"{OCIS_URL}/remote.php/webdav/test.wopitest",
         "-D", headers_file],
        check=True,
    )

    # Extract Oc-Fileid from response headers
    headers_text = Path(headers_file).read_text()
    print("--- PUT headers ---", flush=True)
    print(headers_text[:500], flush=True)
    m = re.search(r"Oc-Fileid:\s*(\S+)", headers_text, re.IGNORECASE)
    if not m:
        print("ERROR: Oc-Fileid not found in PUT response headers", file=sys.stderr)
        sys.exit(1)
    file_id = m.group(1).strip()
    print(f"FILE_ID={file_id}", flush=True)

    # POST to app/open to get WOPI access token and wopi src
    url = f"{OCIS_URL}/app/open?app_name=FakeOffice&file_id={urllib.parse.quote(file_id, safe='')}"
    r = subprocess.run(
        ["curl", "-sk", "-u", "admin:admin", "-X", "POST",
         "--fail", "--retry-connrefused", "--retry", "7", "--retry-all-errors", url],
        capture_output=True, text=True, check=True,
    )
    open_json = json.loads(r.stdout)
    print(f"open.json: {r.stdout[:800]}", flush=True)

    access_token = open_json["form_parameters"]["access_token"]
    access_token_ttl = str(open_json["form_parameters"]["access_token_ttl"])
    app_url = open_json.get("app_url", "")

    # Construct wopi_src: drone extracts file ID from app_url after 'files%2F',
    # then prepends http://wopi-fakeoffice:9300/wopi/files/ — we use bridge_ip instead.
    wopi_base = f"http://{bridge_ip}:9300/wopi/files/"
    if "files%2F" in app_url:
        file_id_encoded = app_url.split("files%2F")[-1].strip().strip('"')
    elif "files/" in app_url:
        file_id_encoded = app_url.split("files/")[-1].strip().strip('"')
    else:
        file_id_encoded = urllib.parse.quote(file_id, safe="")
    wopi_src = wopi_base + file_id_encoded
    print(f"WOPI_SRC={wopi_src}", flush=True)

    return access_token, access_token_ttl, wopi_src


def run_validator(group: str, token: str, wopi_src: str, ttl: str,
                  secure: bool = False) -> int:
    print(f"\nRunning testgroup [{group}] secure={secure}", flush=True)
    cmd = [
        "docker", "run", "--rm",
        "--workdir", "/app",
        "--entrypoint", "/app/Microsoft.Office.WopiValidator",
        VALIDATOR_IMAGE,
    ]
    if secure:
        cmd.append("-s")
    cmd += ["-t", token, "-w", wopi_src, "-l", ttl, "--testgroup", group]
    return subprocess.run(cmd).returncode


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--type", choices=["builtin", "cs3"], required=True,
                        help="WOPI server type: builtin (collaboration service) or cs3 (cs3org/wopiserver)")
    args = parser.parse_args()
    wopi_type = args.type

    repo_root = Path(__file__).resolve().parents[2]
    ocis_bin = repo_root / "ocis/bin/ocis"
    ocis_config_dir = Path.home() / ".ocis/config"

    subprocess.run(["make", "-C", str(repo_root / "ocis"), "build"], check=True)

    bridge_ip = get_docker_bridge_ip()
    print(f"Docker bridge IP: {bridge_ip}", flush=True)

    procs = []
    containers = []

    def cleanup(*_):
        for p in procs:
            try:
                p.terminate()
            except Exception:
                pass
        for name in containers:
            subprocess.run(["docker", "rm", "-f", name], capture_output=True)

    signal.signal(signal.SIGTERM, cleanup)
    signal.signal(signal.SIGINT, cleanup)

    try:
        # --- fakeoffice: serves hosting-discovery.xml on :8080 ---
        # Mirrors drone fakeOffice() — owncloudci/alpine running serve-hosting-discovery.sh.
        # Repo is mounted at /drone/src (the path the script uses).
        containers.append("wopi-fakeoffice-fake")
        subprocess.run(["docker", "rm", "-f", "wopi-fakeoffice-fake"], capture_output=True)
        subprocess.run([
            "docker", "run", "-d", "--name", "wopi-fakeoffice-fake",
            "-p", "8080:8080",
            "-v", f"{repo_root}:/drone/src",
            FAKEOFFICE_IMAGE,
            "sh", "/drone/src/tests/config/drone/serve-hosting-discovery.sh",
        ], check=True)

        wait_for(lambda: tcp_reachable(bridge_ip, 8080), 60, "fakeoffice:8080")
        print("fakeoffice ready.", flush=True)

        # --- Init and start OCIS ---
        ocis_public_url = f"https://{bridge_ip}:9200"
        server_env = {**os.environ}
        server_env.update(base_server_env(
            repo_root, str(ocis_config_dir), ocis_public_url, bridge_ip, wopi_type))

        subprocess.run(
            [str(ocis_bin), "init", "--insecure", "true"],
            env=server_env, check=True,
        )
        shutil.copy(
            repo_root / "tests/config/drone/app-registry.yaml",
            ocis_config_dir / "app-registry.yaml",
        )

        print("Starting ocis...", flush=True)
        ocis_proc = subprocess.Popen([str(ocis_bin), "server"], env=server_env)
        procs.append(ocis_proc)

        wait_for(lambda: ocis_healthy(OCIS_URL), 300, "ocis")
        print("ocis ready.", flush=True)

        # --- Wait for fakeoffice discovery endpoint before starting WOPI service ---
        # ocis collaboration server calls GetAppURLs synchronously at startup;
        # if /hosting/discovery returns non-200, the process exits immediately.
        wait_for(lambda: wopi_discovery_ready("http://127.0.0.1:8080/hosting/discovery"),
                 300, "fakeoffice /hosting/discovery")
        print("fakeoffice discovery ready.", flush=True)

        # --- Start wopi server (after OCIS is healthy so NATS/gRPC are up) ---
        if wopi_type == "builtin":
            # Run 'ocis collaboration server' as a host process.
            # Mirrors drone wopiCollaborationService("fakeoffice") → startOcisService("collaboration").
            collab_env = {**os.environ}
            collab_env.update(collab_service_env(bridge_ip, str(ocis_config_dir)))
            print("Starting collaboration service...", flush=True)
            collab_proc = subprocess.Popen(
                [str(ocis_bin), "collaboration", "server"],
                env=collab_env,
            )
            procs.append(collab_proc)
        else:
            # cs3: patch wopiserver.conf (replace container hostname with bridge_ip),
            # then run cs3org/wopiserver as a Docker container.
            conf_text = (repo_root / "tests/config/drone/wopiserver.conf").read_text()
            conf_text = conf_text.replace("ocis-server", bridge_ip)
            conf_tmp = Path("/tmp/wopiserver-patched.conf")
            conf_tmp.write_text(conf_text)
            secret_tmp = Path("/tmp/wopisecret")
            secret_tmp.write_text("123\n")

            containers.append("wopi-cs3server")
            subprocess.run(["docker", "rm", "-f", "wopi-cs3server"], capture_output=True)
            subprocess.run([
                "docker", "run", "-d", "--name", "wopi-cs3server",
                "-p", "9300:9300",
                "-v", f"{conf_tmp}:/etc/wopi/wopiserver.conf",
                "-v", f"{secret_tmp}:/etc/wopi/wopisecret",
                "--entrypoint", "/app/wopiserver.py",
                CS3_WOPI_IMAGE,
            ], check=True)

        wait_for(lambda: tcp_reachable(bridge_ip, 9300), 120, "wopi-fakeoffice:9300")
        print("wopi server ready.", flush=True)

        # --- prepare-test-file: upload file, get WOPI credentials ---
        access_token, ttl, wopi_src = prepare_test_file(bridge_ip)

        # --- Run validator for each testgroup ---
        failed = []
        for group in SHARED_TESTGROUPS:
            rc = run_validator(group, access_token, wopi_src, ttl, secure=False)
            if rc != 0:
                failed.append(group)

        if wopi_type == "builtin":
            for group in BUILTIN_ONLY_TESTGROUPS:
                rc = run_validator(group, access_token, wopi_src, ttl, secure=True)
                if rc != 0:
                    failed.append(group)

        if failed:
            print(f"\nFailed testgroups: {', '.join(failed)}", file=sys.stderr)
            return 1
        print("\nAll WOPI validator tests passed.")
        return 0

    finally:
        cleanup()


if __name__ == "__main__":
    sys.exit(main())
