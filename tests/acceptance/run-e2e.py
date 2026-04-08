#!/usr/bin/env python3
"""
Run Playwright e2e tests against a local OCIS instance.

Usage: E2E_ARGS='--run-part 1' python3 tests/acceptance/run-e2e.py
Optional: TIKA_NEEDED=true, KEYCLOAK_NEEDED=true
"""

import json
import os
import sys
import shlex
import subprocess
import signal
import time
import shutil
from pathlib import Path

WEB_REPO = "https://github.com/owncloud/web.git"


def wait_for(condition_fn, timeout: int, label: str) -> None:
    deadline = time.time() + timeout
    while not condition_fn():
        if time.time() > deadline:
            print(f"Timeout waiting for {label}", file=sys.stderr)
            sys.exit(1)
        time.sleep(1)


def ocis_healthy(ocis_url: str, use_basic_auth: bool = True) -> bool:
    if use_basic_auth:
        cmd = ["curl", "-sk", "-uadmin:admin",
               f"{ocis_url}/graph/v1.0/users/admin",
               "-w", "%{http_code}", "-o", "/dev/null"]
    else:
        # Keycloak mode: no local admin user, check unauthenticated endpoint
        cmd = ["curl", "-sk",
               f"{ocis_url}/.well-known/openid-configuration",
               "-w", "%{http_code}", "-o", "/dev/null"]
    r = subprocess.run(cmd, capture_output=True, text=True)
    return r.stdout.strip() == "200"


def run(cmd: list, env: dict = None, check: bool = True, cwd=None):
    e = {**os.environ, **(env or {})}
    return subprocess.run(cmd, env=e, check=check, cwd=cwd)


def main() -> int:
    e2e_args = os.environ.get("E2E_ARGS", "").strip()
    if not e2e_args:
        print("E2E_ARGS is required, e.g. E2E_ARGS='--run-part 1' python3 run-e2e.py",
              file=sys.stderr)
        return 1

    tika_needed = os.environ.get("TIKA_NEEDED", "").lower() == "true"
    keycloak_needed = os.environ.get("KEYCLOAK_NEEDED", "").lower() == "true"

    repo_root = Path(__file__).resolve().parents[2]
    ocis_bin = repo_root / "ocis/bin/ocis"
    wrapper_bin = repo_root / "tests/ociswrapper/bin/ociswrapper"
    ocis_url = "https://127.0.0.1:9200"
    ocis_config_dir = Path.home() / ".ocis/config"
    web_dir = repo_root / "webTestRunner"

    # build ocis + ociswrapper only if not already provided (e.g. via artifact)
    if not ocis_bin.exists():
        run(["make", "-C", str(repo_root / "ocis"), "build"])
    if not wrapper_bin.exists():
        run(["make", "-C", str(repo_root / "tests/ociswrapper"), "build"],
            env={"GOWORK": "off"})

    # clone + install web only if not already provided (e.g. via artifact)
    if not web_dir.exists():
        drone_env = {}
        drone_env_file = repo_root / ".drone.env"
        if drone_env_file.exists():
            for line in drone_env_file.read_text().splitlines():
                if "=" in line and not line.startswith("#"):
                    k, v = line.split("=", 1)
                    drone_env[k.strip()] = v.strip()
        web_branch = drone_env.get("WEB_BRANCH", "master")
        web_commitid = drone_env.get("WEB_COMMITID", "")
        run(["git", "clone", "-b", web_branch, "--single-branch", "--no-tags",
             WEB_REPO, str(web_dir)])
        if web_commitid:
            subprocess.run(["git", "checkout", web_commitid], cwd=web_dir, check=True)

    if not (web_dir / "node_modules").exists():
        pkg_manager = json.loads((web_dir / "package.json").read_text()).get("packageManager", "pnpm")
        run(["npm", "install", "--silent", "--global", "--force", pkg_manager])
        subprocess.run(["pnpm", "config", "set", "store-dir", "./.pnpm-store"],
                       cwd=web_dir, check=True)
        subprocess.run(["pnpm", "install"], cwd=web_dir, check=True)
        subprocess.run(["pnpm", "exec", "playwright", "install", "--with-deps", "chromium"],
                       cwd=web_dir, check=True)

    # init ocis
    run([str(ocis_bin), "init", "--insecure", "true"])
    shutil.copy(
        repo_root / "tests/config/drone/app-registry.yaml",
        ocis_config_dir / "app-registry.yaml",
    )

    # Trust the self-signed cert system-wide so Chromium and node-fetch
    # don't produce TLS handshake errors that load the server
    ocis_cert = Path.home() / ".ocis/proxy/server.crt"
    if ocis_cert.exists():
        subprocess.run(
            ["sudo", "cp", str(ocis_cert), "/usr/local/share/ca-certificates/ocis.crt"],
            check=True,
        )
        subprocess.run(["sudo", "update-ca-certificates"], check=True)

    # Patch web UI config: replace Drone Docker service name with our URL
    drone_web_cfg = json.loads(
        (repo_root / "tests/config/drone/ocis-config.json").read_text()
    )

    def _patch_urls(obj, old, new):
        if isinstance(obj, dict):
            return {k: _patch_urls(v, old, new) for k, v in obj.items()}
        if isinstance(obj, list):
            return [_patch_urls(v, old, new) for v in obj]
        if isinstance(obj, str):
            return obj.replace(old, new)
        return obj

    gha_web_cfg = _patch_urls(drone_web_cfg, "https://ocis-server:9200", ocis_url)
    gha_web_cfg_path = ocis_config_dir / "web-ui-config.json"
    gha_web_cfg_path.write_text(json.dumps(gha_web_cfg, indent=2))

    server_env = {
        **os.environ,
        # core — matches drone ocisServer()
        "OCIS_URL": ocis_url,
        "OCIS_CONFIG_DIR": str(ocis_config_dir),
        "STORAGE_USERS_DRIVER": "ocis",
        "PROXY_ENABLE_BASIC_AUTH": "true",

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
        "OCIS_TRANSLATION_PATH": str(repo_root / "tests/config/translations"),
        "WEB_UI_CONFIG_FILE": str(gha_web_cfg_path),
        "THUMBNAILS_TXT_FONTMAP_FILE": str(repo_root / "tests/config/drone/fontsMap.json"),
        # extra_server_environment — matches drone e2eTestPipeline()
        "OCIS_PASSWORD_POLICY_BANNED_PASSWORDS_LIST": str(repo_root / "tests/config/drone/banned-password-list.txt"),
        "GRAPH_AVAILABLE_ROLES": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5,a8d5fe5e-96e3-418d-825b-534dbdf22b99,fb6c3e19-e378-47e5-b277-9732f9de6e21,58c63c02-1d89-4572-916a-870abc5a1b7d,2d00ce52-1fc2-4dbc-8b95-a73b73395f5a,1c996275-f1c9-4e71-abdf-a42f6495e960,312c0871-5ef7-4b3a-85b6-0e4074c64049,aa97fe03-7980-45ac-9e50-b325749fd7e6,63e64e19-8d43-42ec-a738-2b6af2610efa",
        "FRONTEND_CONFIGURABLE_NOTIFICATIONS": "true",
        # debug addresses
        "ACTIVITYLOG_DEBUG_ADDR": "0.0.0.0:9197",
        "APP_PROVIDER_DEBUG_ADDR": "0.0.0.0:9165",
        "APP_REGISTRY_DEBUG_ADDR": "0.0.0.0:9243",
        "AUTH_BASIC_DEBUG_ADDR": "0.0.0.0:9147",
        "AUTH_MACHINE_DEBUG_ADDR": "0.0.0.0:9167",
        "AUTH_SERVICE_DEBUG_ADDR": "0.0.0.0:9198",
        "CLIENTLOG_DEBUG_ADDR": "0.0.0.0:9260",
        "EVENTHISTORY_DEBUG_ADDR": "0.0.0.0:9270",
        "FRONTEND_DEBUG_ADDR": "0.0.0.0:9141",
        "GATEWAY_DEBUG_ADDR": "0.0.0.0:9143",
        "GRAPH_DEBUG_ADDR": "0.0.0.0:9124",
        "GROUPS_DEBUG_ADDR": "0.0.0.0:9161",
        "IDM_DEBUG_ADDR": "0.0.0.0:9239",
        "IDP_DEBUG_ADDR": "0.0.0.0:9134",
        "INVITATIONS_DEBUG_ADDR": "0.0.0.0:9269",
        "NATS_DEBUG_ADDR": "0.0.0.0:9234",
        "OCDAV_DEBUG_ADDR": "0.0.0.0:9163",
        "OCM_DEBUG_ADDR": "0.0.0.0:9281",
        "OCS_DEBUG_ADDR": "0.0.0.0:9114",
        "POSTPROCESSING_DEBUG_ADDR": "0.0.0.0:9255",
        "PROXY_DEBUG_ADDR": "0.0.0.0:9205",
        "SEARCH_DEBUG_ADDR": "0.0.0.0:9224",
        "SETTINGS_DEBUG_ADDR": "0.0.0.0:9194",
        "SHARING_DEBUG_ADDR": "0.0.0.0:9151",
        "SSE_DEBUG_ADDR": "0.0.0.0:9139",
        "STORAGE_PUBLICLINK_DEBUG_ADDR": "0.0.0.0:9179",
        "STORAGE_SHARES_DEBUG_ADDR": "0.0.0.0:9156",
        "STORAGE_SYSTEM_DEBUG_ADDR": "0.0.0.0:9217",
        "STORAGE_USERS_DEBUG_ADDR": "0.0.0.0:9159",
        "THUMBNAILS_DEBUG_ADDR": "0.0.0.0:9189",
        "USERLOG_DEBUG_ADDR": "0.0.0.0:9214",
        "USERS_DEBUG_ADDR": "0.0.0.0:9145",
        "WEB_DEBUG_ADDR": "0.0.0.0:9104",
        "WEBDAV_DEBUG_ADDR": "0.0.0.0:9119",
        "WEBFINGER_DEBUG_ADDR": "0.0.0.0:9279",
    }

    if tika_needed:
        server_env.update({
            "FRONTEND_FULL_TEXT_SEARCH_ENABLED": "true",
            "SEARCH_EXTRACTOR_TYPE": "tika",
            "SEARCH_EXTRACTOR_TIKA_TIKA_URL": "http://localhost:9998",
            "SEARCH_EXTRACTOR_CS3SOURCE_INSECURE": "true",
        })

    if keycloak_needed:
        server_env.update({
            "OCIS_EXCLUDE_RUN_SERVICES": "idp",
            "PROXY_AUTOPROVISION_ACCOUNTS": "true",
            "PROXY_ROLE_ASSIGNMENT_DRIVER": "oidc",
            "OCIS_OIDC_ISSUER": "https://localhost:8443/realms/oCIS",
            "PROXY_OIDC_REWRITE_WELLKNOWN": "true",
            "WEB_OIDC_CLIENT_ID": "web",
            "PROXY_USER_OIDC_CLAIM": "preferred_username",
            "PROXY_USER_CS3_CLAIM": "username",
            "OCIS_ADMIN_USER_ID": "",
            "GRAPH_ASSIGN_DEFAULT_USER_ROLE": "false",
            "GRAPH_USERNAME_MATCH": "none",
            "PROXY_CSP_CONFIG_FILE_LOCATION": str(repo_root / "tests/config/drone/csp.yaml"),
            "KEYCLOAK_DOMAIN": "localhost:8443",
            "IDM_CREATE_DEMO_USERS": "false",
        })

    procs = []

    print("Starting ocis...")
    if keycloak_needed:
        # external IdP: run ocis server directly (no ociswrapper)
        ocis_proc = subprocess.Popen(
            [str(ocis_bin), "server"],
            env=server_env,
        )
    else:
        ocis_proc = subprocess.Popen(
            [str(wrapper_bin), "serve",
             "--bin", str(ocis_bin),
             "--url", ocis_url,
             "--admin-username", "admin",
             "--admin-password", "admin"],
            env=server_env,
        )
    procs.append(ocis_proc)

    def cleanup(*_):
        for p in procs:
            try:
                p.terminate()
            except Exception:
                pass

    signal.signal(signal.SIGTERM, cleanup)
    signal.signal(signal.SIGINT, cleanup)

    try:
        wait_for(lambda: ocis_healthy(ocis_url, use_basic_auth=not keycloak_needed), 300, "ocis")
        print("ocis ready.")

        playwright_env = {
            **os.environ,
            "BASE_URL_OCIS": ocis_url,
            "HEADLESS": "true",
            "RETRY": "3",
            "SKIP_A11Y_TESTS": "true",
            "REPORT_TRACING": "true",
            "NODE_EXTRA_CA_CERTS": str(ocis_cert),
            "BROWSER": "chromium",
        }

        if keycloak_needed:
            playwright_env.update({
                "KEYCLOAK": "true",
                "KEYCLOAK_HOST": "localhost:8443",
            })

        print(f"Running e2e: {e2e_args}")
        result = subprocess.run(
            ["bash", "run-e2e.sh"] + shlex.split(e2e_args),
            cwd=web_dir / "tests/e2e",
            env=playwright_env,
        )
        return result.returncode

    finally:
        cleanup()


if __name__ == "__main__":
    sys.exit(main())
