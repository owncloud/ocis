#!/usr/bin/env python3
"""
Run ocis acceptance tests locally and in GitHub Actions CI.

Config sourced from .drone.star localApiTests — single source of truth.
Usage: BEHAT_SUITES=apiGraph python3 tests/acceptance/run-github.py
"""

import json
import os
import re
import sys
import subprocess
import signal
import time
import tempfile
import shutil
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path

# ---------------------------------------------------------------------------
# Config sourced from .drone.star
# NOTE: EMAIL_SMTP_HOST is "email" (container name) in drone, "localhost" here
# ---------------------------------------------------------------------------

EMAIL_SMTP_HOST = "localhost"
EMAIL_SMTP_PORT = "1025"
EMAIL_PORT = "8025"
EMAIL_SMTP_SENDER = "ownCloud <noreply@example.com>"

LOCAL_API_TESTS = {
    "contractAndLock": {
        "suites": ["apiContract", "apiLocks"],
    },
    "settingsAndNotification": {
        "suites": ["apiSettings", "apiNotification", "apiCors"],
        "emailNeeded": True,
        "extraEnvironment": {
            "EMAIL_HOST": EMAIL_SMTP_HOST,
            "EMAIL_PORT": EMAIL_PORT,
        },
        "extraServerEnvironment": {
            "OCIS_ADD_RUN_SERVICES": "notifications",
            "NOTIFICATIONS_SMTP_HOST": EMAIL_SMTP_HOST,
            "NOTIFICATIONS_SMTP_PORT": EMAIL_SMTP_PORT,
            "NOTIFICATIONS_SMTP_INSECURE": "true",
            "NOTIFICATIONS_SMTP_SENDER": EMAIL_SMTP_SENDER,
            "NOTIFICATIONS_DEBUG_ADDR": "0.0.0.0:9174",
        },
    },
    "graphUser": {
        "suites": ["apiGraphUser"],
    },
    "spaces": {
        "suites": ["apiSpaces"],
    },
    "spacesShares": {
        "suites": ["apiSpacesShares"],
    },
    "davOperations": {
        "suites": [
            "apiSpacesDavOperation", "apiDownloads", "apiAsyncUpload",
            "apiDepthInfinity", "apiArchiver", "apiActivities",
        ],
    },
    "groupAndSearch1": {
        "suites": ["apiSearch1", "apiGraph", "apiGraphGroup"],
    },
    "search2": {
        "suites": ["apiSearch2", "apiSearchContent"],
        "tikaNeeded": True,
        "extraServerEnvironment": {
            "FRONTEND_FULL_TEXT_SEARCH_ENABLED": "true",
            "SEARCH_EXTRACTOR_TYPE": "tika",
            "SEARCH_EXTRACTOR_TIKA_TIKA_URL": "http://localhost:9998",
            "SEARCH_EXTRACTOR_CS3SOURCE_INSECURE": "true",
        },
    },
    "sharingNg1": {
        "suites": ["apiSharingNgShares", "apiReshare", "apiSharingNgPermissions"],
    },
    "sharingNgAdditionalShareRole": {
        "suites": ["apiSharingNgAdditionalShareRole"],
    },
    "sharingNgShareInvitation": {
        "suites": ["apiSharingNgDriveInvitation", "apiSharingNgItemInvitation"],
    },
    "sharingNgLinkShare": {
        "suites": [
            "apiSharingNgDriveLinkShare", "apiSharingNgItemLinkShare",
            "apiSharingNgLinkShareManagement",
        ],
    },
    "antivirus": {
        "suites": ["apiAntivirus"],
        "antivirusNeeded": True,
        "extraServerEnvironment": {
            "ANTIVIRUS_SCANNER_TYPE": "clamav",
            "ANTIVIRUS_CLAMAV_SOCKET": "tcp://clamav:3310",
            "POSTPROCESSING_STEPS": "virusscan",
            "OCIS_ADD_RUN_SERVICES": "antivirus",
            "ANTIVIRUS_DEBUG_ADDR": "0.0.0.0:9277",
        },
    },
    "ocm": {
        "suites": ["apiOcm"],
        "emailNeeded": True,
        "federationServer": True,
        "extraEnvironment": {
            "EMAIL_HOST": EMAIL_SMTP_HOST,
            "EMAIL_PORT": EMAIL_PORT,
        },
        "extraServerEnvironment": {
            "OCIS_ADD_RUN_SERVICES": "ocm,notifications",
            "OCIS_ENABLE_OCM": "true",
            "OCM_OCM_INVITE_MANAGER_INSECURE": "true",
            "OCM_OCM_SHARE_PROVIDER_INSECURE": "true",
            "OCM_OCM_STORAGE_PROVIDER_INSECURE": "true",
            "OCM_OCM_PROVIDER_AUTHORIZER_PROVIDERS_FILE": "",  # set at runtime
            "NOTIFICATIONS_SMTP_HOST": EMAIL_SMTP_HOST,
            "NOTIFICATIONS_SMTP_PORT": EMAIL_SMTP_PORT,
            "NOTIFICATIONS_SMTP_INSECURE": "true",
            "NOTIFICATIONS_SMTP_SENDER": EMAIL_SMTP_SENDER,
            "NOTIFICATIONS_DEBUG_ADDR": "0.0.0.0:9174",
        },
    },
    "authApp": {
        "suites": ["apiAuthApp"],
        "extraServerEnvironment": {
            "OCIS_ADD_RUN_SERVICES": "auth-app",
            "PROXY_ENABLE_APP_AUTH": "true",
        },
    },
    "wopi": {
        "suites": ["apiCollaboration"],
        "collaborationServiceNeeded": True,
        "extraServerEnvironment": {
            "GATEWAY_GRPC_ADDR": "0.0.0.0:9142",
        },
    },
    "cliCommands": {
        "suites": ["cliCommands", "apiServiceAvailability"],
        "antivirusNeeded": True,
        "emailNeeded": True,
        "extraEnvironment": {
            "EMAIL_HOST": EMAIL_SMTP_HOST,
            "EMAIL_PORT": EMAIL_PORT,
        },
        "extraServerEnvironment": {
            "NOTIFICATIONS_SMTP_HOST": EMAIL_SMTP_HOST,
            "NOTIFICATIONS_SMTP_PORT": EMAIL_SMTP_PORT,
            "NOTIFICATIONS_SMTP_INSECURE": "true",
            "NOTIFICATIONS_SMTP_SENDER": EMAIL_SMTP_SENDER,
            "NOTIFICATIONS_DEBUG_ADDR": "0.0.0.0:9174",
            "ANTIVIRUS_SCANNER_TYPE": "clamav",
            "ANTIVIRUS_CLAMAV_SOCKET": "tcp://clamav:3310",
            "ANTIVIRUS_DEBUG_ADDR": "0.0.0.0:9277",
            "OCIS_ADD_RUN_SERVICES": "antivirus,notifications",
        },
    },
}

# reverse lookup: suite → group config
_SUITE_TO_CONFIG: dict = {}
for _cfg in LOCAL_API_TESTS.values():
    for _s in _cfg.get("suites", []):
        _SUITE_TO_CONFIG[_s] = _cfg


# GitHub Actions uses --network host: all wopi services share one network namespace.
# Drone gives each service its own container → all can use 9300/9301/9304.
# Assign distinct ports here to avoid collisions.
_WOPI_PORTS = {
    "collabora":  {"grpc": 9301, "http": 9300, "debug": 9304},
    "onlyoffice": {"grpc": 9311, "http": 9310, "debug": 9314},
    "fakeoffice": {"grpc": 9321, "http": 9320, "debug": 9324},
}


def merged_config(suites: list) -> dict:
    """Union config for all requested suites."""
    merged = {
        "emailNeeded": False,
        "antivirusNeeded": False,
        "tikaNeeded": False,
        "federationServer": False,
        "collaborationServiceNeeded": False,
        "extraServerEnvironment": {},
        "extraEnvironment": {},
    }
    for suite in suites:
        cfg = _SUITE_TO_CONFIG.get(suite, {})
        for flag in ("emailNeeded", "antivirusNeeded", "tikaNeeded",
                     "federationServer", "collaborationServiceNeeded"):
            if cfg.get(flag):
                merged[flag] = True
        merged["extraServerEnvironment"].update(cfg.get("extraServerEnvironment", {}))
        merged["extraEnvironment"].update(cfg.get("extraEnvironment", {}))
    return merged


def _apply_port_offset(env: dict, offset: int) -> dict:
    """Offset OCIS service ports (9100–9399) embedded in env values like '0.0.0.0:9174'."""
    if offset == 0:
        return env
    result = {}
    for k, v in env.items():
        m = re.match(r'^(.*:)(\d{4,5})$', str(v))
        if m:
            port = int(m.group(2))
            if 9100 <= port < 9500:  # covers all OCIS ports incl. store (9460/9464)
                v = m.group(1) + str(port + offset)
        result[k] = v
    return result


def _ocis_slot_dirs(slot: int):
    """Return (config_dir, data_dir) for the given isolation slot."""
    if slot == 0:
        return Path.home() / ".ocis" / "config", Path.home() / ".ocis"
    return (Path.home() / f".ocis-slot-{slot}" / "config",
            Path.home() / f".ocis-slot-{slot}")


def base_server_env(repo_root: Path, ocis_url: str, ocis_config_dir: str,
                    port_offset: int = 0) -> dict:
    """Base ocis server environment matching drone ocisServer() function."""
    def p(base: int) -> str:
        return str(base + port_offset)

    return {
        "OCIS_URL": ocis_url,
        "OCIS_CONFIG_DIR": ocis_config_dir,
        "STORAGE_USERS_DRIVER": "ocis",
        "PROXY_ENABLE_BASIC_AUTH": "true",
        "OCIS_LOG_LEVEL": "error",
        "IDM_CREATE_DEMO_USERS": "true",
        "IDM_ADMIN_PASSWORD": "admin",
        "FRONTEND_SEARCH_MIN_LENGTH": "2",
        "OCIS_ASYNC_UPLOADS": "true",
        "OCIS_EVENTS_ENABLE_TLS": "false",
        "NATS_NATS_HOST": "0.0.0.0",
        "NATS_NATS_PORT": p(9233),
        "MICRO_REGISTRY_ADDRESS": f"127.0.0.1:{p(9233)}",
        "OCIS_RUNTIME_PORT": p(9250),
        "OCIS_JWT_SECRET": "some-ocis-jwt-secret",
        "EVENTHISTORY_STORE": "memory",
        "OCIS_TRANSLATION_PATH": str(repo_root / "tests/config/translations"),
        "WEB_UI_CONFIG_FILE": str(repo_root / "tests/config/drone/ocis-config.json"),
        "THUMBNAILS_TXT_FONTMAP_FILE": str(repo_root / "tests/config/drone/fontsMap.json"),
        # default tika off (overridden by search2 extraServerEnvironment)
        "SEARCH_EXTRACTOR_TYPE": "basic",
        "FRONTEND_FULL_TEXT_SEARCH_ENABLED": "false",
        # OCIS events endpoint (NATS) — each slot needs its own NATS port
        "OCIS_EVENTS_ENDPOINT": f"127.0.0.1:{p(9233)}",
        # LDAP (IDM) — each slot needs its own LDAP listener
        "OCIS_LDAP_URI":                 f"ldaps://localhost:{p(9235)}",
        "IDM_LDAPS_ADDR":                f"0.0.0.0:{p(9235)}",
        # proxy HTTP — must be set explicitly so each slot binds to its own port
        "PROXY_HTTP_ADDR":               f"0.0.0.0:{p(9200)}",
        # gRPC listen addresses — ALL services; same port-block reasoning as debug addrs
        "APP_PROVIDER_GRPC_ADDR":        f"0.0.0.0:{p(9164)}",
        "APP_REGISTRY_GRPC_ADDR":        f"0.0.0.0:{p(9242)}",
        "AUTH_BASIC_GRPC_ADDR":          f"0.0.0.0:{p(9146)}",
        "AUTH_MACHINE_GRPC_ADDR":        f"0.0.0.0:{p(9166)}",
        "AUTH_SERVICE_GRPC_ADDR":        f"0.0.0.0:{p(9199)}",
        "EVENTHISTORY_GRPC_ADDR":        f"0.0.0.0:{p(9274)}",
        "GATEWAY_GRPC_ADDR":             f"0.0.0.0:{p(9142)}",
        "GROUPS_GRPC_ADDR":              f"0.0.0.0:{p(9160)}",
        "SEARCH_GRPC_ADDR":              f"0.0.0.0:{p(9220)}",
        "SETTINGS_GRPC_ADDR":            f"0.0.0.0:{p(9185)}",
        "SHARING_GRPC_ADDR":             f"0.0.0.0:{p(9150)}",
        "STORAGE_PUBLICLINK_GRPC_ADDR":  f"0.0.0.0:{p(9178)}",
        "STORAGE_SHARES_GRPC_ADDR":      f"0.0.0.0:{p(9154)}",
        "STORAGE_SYSTEM_GRPC_ADDR":      f"0.0.0.0:{p(9215)}",
        "STORAGE_USERS_GRPC_ADDR":       f"0.0.0.0:{p(9157)}",
        "THUMBNAILS_GRPC_ADDR":          f"0.0.0.0:{p(9191)}",
        "USERS_GRPC_ADDR":               f"0.0.0.0:{p(9144)}",
        "STORE_GRPC_ADDR":               f"0.0.0.0:{p(9460)}",
        # HTTP listen addresses
        "ACTIVITYLOG_HTTP_ADDR":         f"0.0.0.0:{p(9195)}",
        "FRONTEND_HTTP_ADDR":            f"0.0.0.0:{p(9140)}",
        "GRAPH_HTTP_ADDR":               f"0.0.0.0:{p(9120)}",
        "IDP_HTTP_ADDR":                 f"0.0.0.0:{p(9130)}",
        "OCDAV_HTTP_ADDR":               f"0.0.0.0:{p(9162)}",
        "OCS_HTTP_ADDR":                 f"0.0.0.0:{p(9110)}",
        "SETTINGS_HTTP_ADDR":            f"0.0.0.0:{p(9186)}",
        "SSE_HTTP_ADDR":                 f"0.0.0.0:{p(9132)}",
        "STORAGE_SYSTEM_HTTP_ADDR":      f"0.0.0.0:{p(9216)}",
        "STORAGE_USERS_HTTP_ADDR":       f"0.0.0.0:{p(9158)}",
        "THUMBNAILS_HTTP_ADDR":          f"0.0.0.0:{p(9190)}",
        "USERLOG_HTTP_ADDR":             f"0.0.0.0:{p(9211)}",
        "WEB_HTTP_ADDR":                 f"0.0.0.0:{p(9100)}",
        "WEBDAV_HTTP_ADDR":              f"0.0.0.0:{p(9115)}",
        "WEBFINGER_HTTP_ADDR":           f"0.0.0.0:{p(9275)}",
        # data server URLs — must reference the slot-specific HTTP port
        "STORAGE_USERS_DATA_SERVER_URL":  f"http://localhost:{p(9158)}/data",
        "STORAGE_SYSTEM_DATA_SERVER_URL": f"http://localhost:{p(9216)}/data",
        "THUMBNAILS_DATA_ENDPOINT":       f"http://127.0.0.1:{p(9190)}/thumbnails/data",
        # debug addresses
        "ACTIVITYLOG_DEBUG_ADDR":        f"0.0.0.0:{p(9197)}",
        "APP_PROVIDER_DEBUG_ADDR":       f"0.0.0.0:{p(9165)}",
        "APP_REGISTRY_DEBUG_ADDR":       f"0.0.0.0:{p(9243)}",
        "AUTH_BASIC_DEBUG_ADDR":         f"0.0.0.0:{p(9147)}",
        "AUTH_MACHINE_DEBUG_ADDR":       f"0.0.0.0:{p(9167)}",
        "AUTH_SERVICE_DEBUG_ADDR":       f"0.0.0.0:{p(9198)}",
        "CLIENTLOG_DEBUG_ADDR":          f"0.0.0.0:{p(9260)}",
        "EVENTHISTORY_DEBUG_ADDR":       f"0.0.0.0:{p(9270)}",
        "FRONTEND_DEBUG_ADDR":           f"0.0.0.0:{p(9141)}",
        "GATEWAY_DEBUG_ADDR":            f"0.0.0.0:{p(9143)}",
        "GRAPH_DEBUG_ADDR":              f"0.0.0.0:{p(9124)}",
        "GROUPS_DEBUG_ADDR":             f"0.0.0.0:{p(9161)}",
        "IDM_DEBUG_ADDR":                f"0.0.0.0:{p(9239)}",
        "IDP_DEBUG_ADDR":                f"0.0.0.0:{p(9134)}",
        "INVITATIONS_DEBUG_ADDR":        f"0.0.0.0:{p(9269)}",
        "NATS_DEBUG_ADDR":               f"0.0.0.0:{p(9234)}",
        "OCDAV_DEBUG_ADDR":              f"0.0.0.0:{p(9163)}",
        "OCM_GRPC_ADDR":                 f"0.0.0.0:{p(9282)}",
        "OCM_HTTP_ADDR":                 f"0.0.0.0:{p(9280)}",
        "OCM_DEBUG_ADDR":                f"0.0.0.0:{p(9281)}",
        "OCS_DEBUG_ADDR":                f"0.0.0.0:{p(9114)}",
        "POSTPROCESSING_DEBUG_ADDR":     f"0.0.0.0:{p(9255)}",
        "PROXY_DEBUG_ADDR":              f"0.0.0.0:{p(9205)}",
        "SEARCH_DEBUG_ADDR":             f"0.0.0.0:{p(9224)}",
        "SETTINGS_DEBUG_ADDR":           f"0.0.0.0:{p(9194)}",
        "SHARING_DEBUG_ADDR":            f"0.0.0.0:{p(9151)}",
        "SSE_DEBUG_ADDR":                f"0.0.0.0:{p(9139)}",
        "STORAGE_PUBLICLINK_DEBUG_ADDR": f"0.0.0.0:{p(9179)}",
        "STORAGE_SHARES_DEBUG_ADDR":     f"0.0.0.0:{p(9156)}",
        "STORAGE_SYSTEM_DEBUG_ADDR":     f"0.0.0.0:{p(9217)}",
        "STORAGE_USERS_DEBUG_ADDR":      f"0.0.0.0:{p(9159)}",
        "STORE_DEBUG_ADDR":              f"0.0.0.0:{p(9464)}",
        "THUMBNAILS_DEBUG_ADDR":         f"0.0.0.0:{p(9189)}",
        "USERLOG_DEBUG_ADDR":            f"0.0.0.0:{p(9214)}",
        "USERS_DEBUG_ADDR":              f"0.0.0.0:{p(9145)}",
        "WEB_DEBUG_ADDR":                f"0.0.0.0:{p(9104)}",
        "WEBDAV_DEBUG_ADDR":             f"0.0.0.0:{p(9119)}",
        "WEBFINGER_DEBUG_ADDR":          f"0.0.0.0:{p(9279)}",
    }


def wait_for(condition_fn, timeout: int, label: str, container: str = None) -> None:
    start = time.time()
    deadline = start + timeout
    last_log = start
    while not condition_fn():
        now = time.time()
        if now > deadline:
            elapsed = int(now - start)
            print(f"Timeout waiting for {label} after {elapsed}s", file=sys.stderr)
            # dump docker diagnostics — use explicit container name if provided
            cname = container or label
            for cmd in (
                ["docker", "ps", "-a", "--filter", f"name={cname}", "--no-trunc"],
                ["docker", "logs", "--tail", "50", cname],
            ):
                r = subprocess.run(cmd, capture_output=True, text=True)
                if r.stdout.strip():
                    print(f"--- {' '.join(cmd)} ---", file=sys.stderr)
                    print(r.stdout, file=sys.stderr)
                if r.stderr.strip():
                    print(r.stderr, file=sys.stderr)
            sys.exit(1)
        if now - last_log >= 30:
            print(f"  Waiting for {label}... {int(now - start)}s")
            last_log = now
        time.sleep(1)


def ocis_healthy(ocis_url: str) -> bool:
    r = subprocess.run(
        ["curl", "-sk", "-uadmin:admin",
         f"{ocis_url}/graph/v1.0/users/admin",
         "-w", "%{http_code}", "-o", "/dev/null"],
        capture_output=True, text=True,
    )
    return r.stdout.strip() == "200"


def mailpit_healthy() -> bool:
    return subprocess.run(
        ["curl", "-sf", "http://localhost:8025/api/v1/messages"],
        capture_output=True,
    ).returncode == 0


def tika_healthy() -> bool:
    return subprocess.run(
        ["curl", "-sf", "http://localhost:9998"],
        capture_output=True,
    ).returncode == 0


def _tcp_ready(host: str, port: int) -> bool:
    """Check if a TCP port is accepting connections."""
    import socket
    try:
        with socket.create_connection((host, port), timeout=2):
            return True
    except (ConnectionRefusedError, OSError):
        return False


def clamav_healthy() -> bool:
    return _tcp_ready("localhost", 3310)



def load_env_file(path: Path) -> dict:
    """Parse a bash-style env file (export KEY=value) into a dict."""
    env = {}
    for line in path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#") or line.startswith("!"):
            continue
        line = line.removeprefix("export ").strip()
        if "=" in line:
            k, v = line.split("=", 1)
            env[k.strip()] = v.strip()
    return env


def run(cmd: list, env: dict = None, check: bool = True):
    e = {**os.environ, **(env or {})}
    return subprocess.run(cmd, env=e, check=check)


def main() -> int:
    behat_suites_raw = os.environ.get("BEHAT_SUITES", "").strip()
    if not behat_suites_raw:
        print("BEHAT_SUITES is required", file=sys.stderr)
        return 1

    suites = [s.strip() for s in behat_suites_raw.split(",") if s.strip()]
    acceptance_test_type = os.environ.get("ACCEPTANCE_TEST_TYPE", "api")

    repo_root = Path(__file__).resolve().parents[2]
    ocis_bin = repo_root / "ocis/bin/ocis"
    wrapper_bin = repo_root / "tests/ociswrapper/bin/ociswrapper"
    ocis_url = "https://localhost:9200"
    ocis_fed_url = "https://localhost:10200"

    cfg = merged_config(suites)
    print(f"Suites: {suites}")
    print(f"Services: email={cfg['emailNeeded']} tika={cfg['tikaNeeded']} "
          f"antivirus={cfg['antivirusNeeded']} federation={cfg['federationServer']} "
          f"wopi={cfg['collaborationServiceNeeded']}")

    # generate IDP web assets (required for IDP service to start; matches drone ci-node-generate)
    run(["make", "-C", str(repo_root / "services/idp"), "ci-node-generate"])
    # download web UI assets (required for robots.txt and other static assets; no pnpm needed)
    run(["make", "-C", str(repo_root / "services/web"), "ci-node-generate"])

    # build (ENABLE_VIPS=true when libvips-dev is installed, matching drone)
    build_env = {}
    if subprocess.run(["pkg-config", "--exists", "vips"],
                      capture_output=True).returncode == 0:
        build_env["ENABLE_VIPS"] = "true"
    run(["make", "-C", str(repo_root / "ocis"), "build"], env=build_env)
    run(["make", "-C", str(repo_root / "tests/ociswrapper"), "build"],
        env={"GOWORK": "off"})

    # php deps
    run(["composer", "install", "--no-progress"],
        env={"COMPOSER_NO_INTERACTION": "1", "COMPOSER_NO_AUDIT": "1"})
    run(["composer", "bin", "behat", "install", "--no-progress"],
        env={"COMPOSER_NO_INTERACTION": "1", "COMPOSER_NO_AUDIT": "1"})

    # optional services
    procs = []

    if cfg["emailNeeded"]:
        print("Starting mailpit...")
        run(["docker", "run", "-d", "--name", "mailpit", "--network", "host",
             "axllent/mailpit:v1.22.3"])
        wait_for(mailpit_healthy, 60, "mailpit")
        print("mailpit ready.")

    if cfg["antivirusNeeded"]:
        print("Starting clamav...")
        run(["docker", "run", "-d", "--name", "clamav", "--network", "host",
             "owncloudci/clamavd"])
        wait_for(clamav_healthy, 300, "clamav")
        print("clamav ready.")
        # override socket: drone uses container DNS "clamav", we use localhost
        cfg["extraServerEnvironment"]["ANTIVIRUS_CLAMAV_SOCKET"] = "tcp://localhost:3310"

    if cfg["tikaNeeded"]:
        print("Starting tika...")
        run(["docker", "run", "-d", "--name", "tika", "--network", "host",
             "apache/tika:3.2.2.0-full"])
        wait_for(tika_healthy, 120, "tika")
        print("tika ready.")

    # OCM federation: rewrite providers.json with localhost URLs
    if cfg["federationServer"]:
        providers_src = repo_root / "tests/config/drone/providers.json"
        providers = json.loads(providers_src.read_text())
        for p in providers:
            # replace container DNS names with localhost
            p["domain"] = p["domain"].replace("ocis-server:9200", "localhost:9200")
            p["domain"] = p["domain"].replace("federation-ocis-server:10200", "localhost:10200")
            for svc in p.get("services", []):
                ep = svc.get("endpoint", {})
                ep["path"] = ep.get("path", "").replace("ocis-server:9200", "localhost:9200")
                ep["path"] = ep.get("path", "").replace("federation-ocis-server:10200", "localhost:10200")
                svc["host"] = svc.get("host", "").replace("ocis-server:9200", "localhost:9200")
                svc["host"] = svc.get("host", "").replace("federation-ocis-server:10200", "localhost:10200")
        providers_tmp = tempfile.NamedTemporaryFile(
            mode="w", suffix=".json", prefix="ocm-providers-", delete=False)
        json.dump(providers, providers_tmp)
        providers_tmp.close()
        cfg["extraServerEnvironment"]["OCM_OCM_PROVIDER_AUTHORIZER_PROVIDERS_FILE"] = providers_tmp.name

    # generate fontsMap.json with correct font path (drone hardcodes /drone/src/...)
    font_path = str(repo_root / "tests/config/drone/NotoSans.ttf")
    fontmap_tmp = tempfile.NamedTemporaryFile(
        mode="w", suffix=".json", prefix="fontsMap-", delete=False)
    json.dump({"defaultFont": font_path}, fontmap_tmp)
    fontmap_tmp.close()

    # Per-slot startup: each suite gets its own isolated OCIS instance.
    # slot 0 uses ~/.ocis (default path, backward compat); slot N uses ~/.ocis-slot-N.
    # Port offset = slot × 300; OCIS debug ports span 9104–9281 (width 177) so offset
    # 300 gives clean separation: slot 0→9100-9399, slot 1→9400-9699, slot 2→9700-9999.
    slot_info = []  # list of (suite, ocis_url, wrapper_url, server_env)
    for slot, suite in enumerate(suites):
        offset = slot * 300
        config_dir, _ = _ocis_slot_dirs(slot)
        config_dir.mkdir(parents=True, exist_ok=True)
        suite_ocis_url = f"https://localhost:{9200 + offset}"
        suite_wrapper_url = f"http://localhost:{5200 + offset}"

        s_env = {**os.environ}
        s_env.update(base_server_env(repo_root, suite_ocis_url, str(config_dir),
                                     port_offset=offset))
        s_env["THUMBNAILS_TXT_FONTMAP_FILE"] = fontmap_tmp.name
        s_env.update(_apply_port_offset(cfg["extraServerEnvironment"], offset))

        # init (idempotent — safe to re-run if config dir already exists)
        run([str(ocis_bin), "init", "--insecure", "true"], env=s_env)
        shutil.copy(
            repo_root / "tests/config/drone/app-registry.yaml",
            config_dir / "app-registry.yaml",
        )

        print(f"Starting ocis slot {slot} (suite: {suite}, port: {9200 + offset}, "
              f"wrapper: {5200 + offset})...")
        wrapper_proc = subprocess.Popen(
            [str(wrapper_bin), "serve",
             "--bin", str(ocis_bin),
             "--url", suite_ocis_url,
             "--admin-username", "admin",
             "--admin-password", "admin",
             "--port", str(5200 + offset)],
            env=s_env,
        )
        procs.append(wrapper_proc)
        slot_info.append((suite, suite_ocis_url, suite_wrapper_url, s_env))

    # convenience: for single-suite callers, keep the old names pointing at slot 0
    ocis_url = slot_info[0][1]

    # start federation ocis server (second instance on port 10200)
    # only used by ocm group (single suite) — always slot 0, no offset conflict
    if cfg["federationServer"]:
        fed_config_dir = Path.home() / ".ocis-federation/config"
        fed_config_dir.mkdir(parents=True, exist_ok=True)
        fed_data_dir = Path.home() / ".ocis-federation"

        fed_env = {**os.environ}
        fed_env.update(base_server_env(repo_root, ocis_fed_url, str(fed_config_dir)))
        fed_env.update(cfg["extraServerEnvironment"])
        # load federation port mappings from canonical env file (single source of truth)
        fed_env.update(load_env_file(repo_root / "tests/config/local/.env-federation"))
        # CI-specific overrides
        fed_env.update({
            "OCIS_URL": ocis_fed_url,
            "OCIS_BASE_DATA_PATH": str(fed_data_dir),
            "OCIS_CONFIG_DIR": str(fed_config_dir),
            "OCIS_RUNTIME_PORT": "10250",
            "MICRO_REGISTRY_ADDRESS": "127.0.0.1:10233",
        })

        # init federation ocis with separate config
        run([str(ocis_bin), "init", "--insecure", "true",
             "--config-path", str(fed_config_dir)])
        shutil.copy(
            repo_root / "tests/config/drone/app-registry.yaml",
            fed_config_dir / "app-registry.yaml",
        )

        print("Starting federation ocis...")
        fed_proc = subprocess.Popen(
            [str(ocis_bin), "server"],
            env=fed_env,
        )
        procs.append(fed_proc)

    # ---------------------------------------------------------------------------
    # Collaboration service helpers — same names/call pattern as drone.star
    # Only deviation: container hostnames → localhost
    # ---------------------------------------------------------------------------
    def fakeOffice():
        # drone: OC_CI_ALPINE container running serve-hosting-discovery.sh with repo at /drone/src
        # use same image so BusyBox nc (not OpenBSD nc) handles FIN correctly on stdin EOF
        run(["docker", "run", "-d", "--name", "fakeoffice", "--network", "host",
             "-v", f"{repo_root}:/drone/src:ro",
             "owncloudci/alpine:latest",
             "sh", "/drone/src/tests/config/drone/serve-hosting-discovery.sh"])
        return []

    def collaboraService():
        # drone commands copy-pasted verbatim
        run(["docker", "run", "-d", "--name", "collabora", "--network", "host",
             "-e", "DONT_GEN_SSL_CERT=set",
             "-e", f"extra_params=--o:ssl.enable=true --o:ssl.termination=true "
                   f"--o:welcome.enable=false --o:net.frame_ancestors=https://localhost:9200",
             "--entrypoint", "/bin/sh",
             "collabora/code:24.04.5.1.1",
             "-c", "\n".join([
                 "set -e",
                 "coolconfig generate-proof-key",
                 "bash /start-collabora-online.sh",
             ])])
        return []

    def onlyofficeService():
        # GitHub runner ships PostgreSQL pre-started on 5432.
        # OnlyOffice supervisord starts its own PostgreSQL on 5432 internally.
        # With --network host both compete for the same port → OnlyOffice DB never
        # starts → docservice stays down → nginx returns 502 forever.
        # Drone avoids this because each service has its own network namespace.
        subprocess.run(["sudo", "systemctl", "stop", "postgresql"],
                       capture_output=True)
        only_office_json = repo_root / "tests/config/drone/only-office.json"
        run(["docker", "run", "-d", "--name", "onlyoffice", "--network", "host",
             "-e", "WOPI_ENABLED=true",
             "-e", "USE_UNAUTHORIZED_STORAGE=true",
             "-v", f"{only_office_json}:/tmp/only-office.json:ro",
             "--entrypoint", "/bin/sh",
             "onlyoffice/documentserver:9.0.0",
             "-c", "\n".join([
                 "set -e",
                 "cp /tmp/only-office.json /etc/onlyoffice/documentserver/local.json",
                 "openssl req -x509 -newkey rsa:4096 -keyout onlyoffice.key -out onlyoffice.crt -sha256 -days 365 -batch -nodes",
                 "mkdir -p /var/www/onlyoffice/Data/certs",
                 "cp onlyoffice.key /var/www/onlyoffice/Data/certs/",
                 "cp onlyoffice.crt /var/www/onlyoffice/Data/certs/",
                 "chmod 400 /var/www/onlyoffice/Data/certs/onlyoffice.key",
                 "/app/ds/run-document-server.sh",
             ])])
        return []

    def wopiCollaborationService(name, ocis_url=ocis_url):
        # drone: startOcisService("collaboration", "wopi-{name}", environment)
        # runs: ocis/bin/ocis-debug collaboration server
        service_name = "wopi-%s" % name
        ports = _WOPI_PORTS[name]
        environment = {
            **os.environ,
            "OCIS_URL": ocis_url,
            "MICRO_REGISTRY": "nats-js-kv",
            "MICRO_REGISTRY_ADDRESS": "localhost:9233",
            "COLLABORATION_LOG_LEVEL": "debug",
            "COLLABORATION_GRPC_ADDR": f"0.0.0.0:{ports['grpc']}",
            "COLLABORATION_HTTP_ADDR": f"0.0.0.0:{ports['http']}",
            "COLLABORATION_DEBUG_ADDR": f"0.0.0.0:{ports['debug']}",
            "COLLABORATION_APP_PROOF_DISABLE": "true",
            "COLLABORATION_APP_INSECURE": "true",
            "COLLABORATION_CS3API_DATAGATEWAY_INSECURE": "true",
            "OCIS_JWT_SECRET": "some-ocis-jwt-secret",
            "COLLABORATION_WOPI_SECRET": "some-wopi-secret",
        }
        if name == "collabora":
            environment["COLLABORATION_APP_NAME"] = "Collabora"
            environment["COLLABORATION_APP_PRODUCT"] = "Collabora"
            environment["COLLABORATION_APP_ADDR"] = "https://localhost:9980"
            environment["COLLABORATION_APP_ICON"] = "https://localhost:9980/favicon.ico"
        elif name == "onlyoffice":
            environment["COLLABORATION_APP_NAME"] = "OnlyOffice"
            environment["COLLABORATION_APP_PRODUCT"] = "OnlyOffice"
            environment["COLLABORATION_APP_ADDR"] = "https://localhost:443"
            environment["COLLABORATION_APP_ICON"] = "https://localhost:443/web-apps/apps/documenteditor/main/resources/img/favicon.ico"
        elif name == "fakeoffice":
            environment["COLLABORATION_APP_NAME"] = "FakeOffice"
            environment["COLLABORATION_APP_PRODUCT"] = "Microsoft"
            environment["COLLABORATION_APP_ADDR"] = "http://localhost:8080"
        environment["COLLABORATION_WOPI_SRC"] = f"http://localhost:{ports['http']}"
        print(f"Starting {service_name}...")
        return [subprocess.Popen([str(ocis_bin), "collaboration", "server"], env=environment)]

    def ocisHealthCheck(name, services=[]):
        # drone: curl healthz + readyz on each service (timeout 300s)
        for service in services:
            host, port = service.rsplit(":", 1)
            for endpoint in ("healthz", "readyz"):
                wait_for(
                    lambda h="localhost", p=int(port), ep=endpoint: subprocess.run(
                        ["curl", "-sf", f"http://{h}:{p}/{ep}"], capture_output=True
                    ).returncode == 0,
                    300, f"{service}/{endpoint}",
                )
        print(f"health-check-{name}: all services healthy.")

    def wopi_discovery_ready(app_url: str) -> bool:
        """Return True once the WOPI app's /hosting/discovery returns HTTP 200."""
        url = app_url.rstrip("/") + "/hosting/discovery"
        r = subprocess.run(
            ["curl", "-sfk", url], capture_output=True
        )
        return r.returncode == 0

    # drone.star non-k8s collaborationServiceNeeded path (lines 1195-1196, 1140-1141, 1179-1183)
    if cfg["collaborationServiceNeeded"]:
        procs += fakeOffice() + collaboraService() + onlyofficeService()
        # Wait for each app's /hosting/discovery to return 200 before starting its
        # collaboration service. GetAppURLs in server.go calls discovery synchronously
        # at startup — non-200 exits the process immediately, healthz never binds.
        wait_for(lambda: wopi_discovery_ready("http://localhost:8080"),  300, "fakeoffice discovery",  container="fakeoffice")
        wait_for(lambda: wopi_discovery_ready("https://localhost:9980"), 300, "collabora discovery",   container="collabora")
        wait_for(lambda: wopi_discovery_ready("https://localhost:443"),  300, "onlyoffice discovery",  container="onlyoffice")
        procs += wopiCollaborationService("fakeoffice") + \
                 wopiCollaborationService("collabora") + \
                 wopiCollaborationService("onlyoffice")
        ocisHealthCheck("wopi", [
            f"localhost:{_WOPI_PORTS['collabora']['debug']}",
            f"localhost:{_WOPI_PORTS['onlyoffice']['debug']}",
            f"localhost:{_WOPI_PORTS['fakeoffice']['debug']}",
        ])

    def cleanup(*_):
        for p in procs:
            try:
                p.terminate()
            except Exception:
                pass

    signal.signal(signal.SIGTERM, cleanup)
    signal.signal(signal.SIGINT, cleanup)

    try:
        for suite, suite_ocis_url, _, _ in slot_info:
            wait_for(lambda u=suite_ocis_url: ocis_healthy(u), 300, f"ocis ({suite})")
            print(f"ocis ready for suite: {suite}.")

        if cfg["federationServer"]:
            wait_for(lambda: ocis_healthy(ocis_fed_url), 300, "federation ocis")
            print("federation ocis ready.")

        # expected failures file
        if acceptance_test_type == "core-api":
            filter_tags = "~@skipOnGraph&&~@skipOnOcis-OCIS-Storage"
            base_failures = repo_root / "tests/acceptance/expected-failures-API-on-OCIS-storage.md"
        else:
            filter_tags = "~@skip&&~@skipOnGraph&&~@skipOnOcis-OCIS-Storage"
            base_failures = repo_root / "tests/acceptance/expected-failures-localAPI-on-OCIS-storage.md"

        ef_override = os.environ.get("EXPECTED_FAILURES_FILE")
        if ef_override:
            p = Path(ef_override)
            base_failures = p if p.is_absolute() else repo_root / p

        # merge expected-failures-without-remotephp.md only when not using remote.php
        # (mirrors drone.star: "" if run_with_remote_php else "cat ...without-remotephp.md >> ...")
        tmp = tempfile.NamedTemporaryFile(mode="w", suffix=".md", delete=False)
        tmp.write(base_failures.read_text())
        if os.environ.get("WITH_REMOTE_PHP", "false").lower() != "true":
            without_rphp = repo_root / "tests/acceptance/expected-failures-without-remotephp.md"
            if without_rphp.exists():
                tmp.write("\n")
                tmp.write(without_rphp.read_text())
        tmp.close()

        # base behat env (suite-specific TEST_SERVER_URL / OCIS_WRAPPER_URL added per slot)
        behat_env_base = {
            **os.environ,
            "TEST_SERVER_FED_URL": ocis_fed_url,
            "ACCEPTANCE_TEST_TYPE": acceptance_test_type,
            "BEHAT_FILTER_TAGS": filter_tags,
            "EXPECTED_FAILURES_FILE": tmp.name,
            "STORAGE_DRIVER": "ocis",
            "UPLOAD_DELETE_WAIT_TIME": "0",
            "EMAIL_HOST": "localhost",
            "EMAIL_PORT": EMAIL_PORT,
            "COLLABORATION_SERVICE_URL": f"http://localhost:{_WOPI_PORTS['fakeoffice']['http']}",
        }
        behat_env_base.update(cfg["extraEnvironment"])

        def _run_suite(suite: str, suite_ocis_url: str, suite_wrapper_url: str) -> tuple:
            """Run behat for one suite; capture output to file; return (suite, rc, log_path)."""
            log_path = Path(f"/tmp/behat-{suite}.log")
            env = {
                **behat_env_base,
                "TEST_SERVER_URL": suite_ocis_url,
                "OCIS_WRAPPER_URL": suite_wrapper_url,
                "BEHAT_SUITES": suite,
            }
            with open(log_path, "w") as lf:
                rc = subprocess.run(
                    ["make", "-C", str(repo_root), "test-acceptance-api"],
                    env=env, stdout=lf, stderr=subprocess.STDOUT,
                ).returncode
            return suite, rc, log_path

        if len(slot_info) == 1:
            # single suite — stream output directly (unchanged behaviour)
            suite, suite_ocis_url, suite_wrapper_url, _ = slot_info[0]
            print(f"Running suite: {suite} (type: {acceptance_test_type})")
            env = {
                **behat_env_base,
                "TEST_SERVER_URL": suite_ocis_url,
                "OCIS_WRAPPER_URL": suite_wrapper_url,
                "BEHAT_SUITES": suite,
            }
            return subprocess.run(
                ["make", "-C", str(repo_root), "test-acceptance-api"], env=env,
            ).returncode

        # multiple suites — run in parallel, print logs serially after all finish
        print(f"Running {len(slot_info)} suites in parallel (type: {acceptance_test_type}): "
              f"{[s for s, *_ in slot_info]}")
        failed = []
        with ThreadPoolExecutor(max_workers=len(slot_info)) as ex:
            futures = {
                ex.submit(_run_suite, suite, suite_ocis_url, suite_wrapper_url): suite
                for suite, suite_ocis_url, suite_wrapper_url, _ in slot_info
            }
            for fut in as_completed(futures):
                suite, rc, log_path = fut.result()
                print(f"\n{'='*60}\n=== Suite: {suite} (rc={rc}) ===\n{'='*60}")
                sys.stdout.write(log_path.read_text())
                sys.stdout.flush()
                if rc != 0:
                    failed.append(suite)

        if failed:
            print(f"\nFailed suites: {', '.join(failed)}", file=sys.stderr)
            return 1
        print(f"\nAll {len(slot_info)} suites passed.")
        return 0

    finally:
        cleanup()


if __name__ == "__main__":
    sys.exit(main())