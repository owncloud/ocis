#!/usr/bin/env python3
"""
Run ocis acceptance tests locally and in GitHub Actions CI.

Config sourced from .drone.star localApiTests — single source of truth.
Usage: BEHAT_SUITES=apiGraph python3 tests/acceptance/run-github.py
"""

import os
import sys
import subprocess
import signal
import time
import tempfile
import shutil
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
            "NOTIFICATIONS_SMTP_HOST": EMAIL_SMTP_HOST,
            "NOTIFICATIONS_SMTP_PORT": EMAIL_SMTP_PORT,
            "NOTIFICATIONS_SMTP_INSECURE": "true",
            "NOTIFICATIONS_SMTP_SENDER": EMAIL_SMTP_SENDER,
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


def base_server_env(repo_root: Path, ocis_url: str, ocis_config_dir: str) -> dict:
    """Base ocis server environment matching drone ocisServer() function."""
    return {
        "OCIS_URL": ocis_url,
        "OCIS_CONFIG_DIR": ocis_config_dir,
        "STORAGE_USERS_DRIVER": "ocis",
        "PROXY_ENABLE_BASIC_AUTH": "true",
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
        "OCIS_TRANSLATION_PATH": str(repo_root / "tests/config/translations"),
        "WEB_UI_CONFIG_FILE": str(repo_root / "tests/config/drone/ocis-config.json"),
        "THUMBNAILS_TXT_FONTMAP_FILE": str(repo_root / "tests/config/drone/fontsMap.json"),
        # default tika off (overridden by search2 extraServerEnvironment)
        "SEARCH_EXTRACTOR_TYPE": "basic",
        "FRONTEND_FULL_TEXT_SEARCH_ENABLED": "false",
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


def clamav_healthy() -> bool:
    """Check ClamAV is ready by attempting a TCP connection to port 3310."""
    import socket
    try:
        with socket.create_connection(("localhost", 3310), timeout=2):
            return True
    except (ConnectionRefusedError, OSError):
        return False


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
    ocis_config_dir = Path.home() / ".ocis/config"

    cfg = merged_config(suites)
    print(f"Suites: {suites}")
    print(f"Services: email={cfg['emailNeeded']} tika={cfg['tikaNeeded']} antivirus={cfg['antivirusNeeded']}")

    # build
    run(["make", "-C", str(repo_root / "ocis"), "build"])
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

    # init ocis
    run([str(ocis_bin), "init", "--insecure", "true"])
    shutil.copy(
        repo_root / "tests/config/drone/app-registry.yaml",
        ocis_config_dir / "app-registry.yaml",
    )

    # assemble ocis server env
    server_env = {**os.environ}
    server_env.update(base_server_env(repo_root, ocis_url, str(ocis_config_dir)))
    server_env.update(cfg["extraServerEnvironment"])

    # start ociswrapper
    print("Starting ocis...")
    wrapper_proc = subprocess.Popen(
        [str(wrapper_bin), "serve",
         "--bin", str(ocis_bin),
         "--url", ocis_url,
         "--admin-username", "admin",
         "--admin-password", "admin"],
        env=server_env,
    )
    procs.append(wrapper_proc)

    def cleanup(*_):
        for p in procs:
            try:
                p.terminate()
            except Exception:
                pass

    signal.signal(signal.SIGTERM, cleanup)
    signal.signal(signal.SIGINT, cleanup)

    try:
        wait_for(lambda: ocis_healthy(ocis_url), 300, "ocis")
        print("ocis ready.")

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

        # merge expected-failures-without-remotephp.md (drone does this)
        tmp = tempfile.NamedTemporaryFile(mode="w", suffix=".md", delete=False)
        tmp.write(base_failures.read_text())
        without_rphp = repo_root / "tests/acceptance/expected-failures-without-remotephp.md"
        if without_rphp.exists():
            tmp.write("\n")
            tmp.write(without_rphp.read_text())
        tmp.close()

        # run tests
        behat_env = {
            **os.environ,
            "TEST_SERVER_URL": ocis_url,
            "OCIS_WRAPPER_URL": "http://localhost:5200",
            "BEHAT_SUITES": behat_suites_raw,
            "ACCEPTANCE_TEST_TYPE": acceptance_test_type,
            "BEHAT_FILTER_TAGS": filter_tags,
            "EXPECTED_FAILURES_FILE": tmp.name,
            "STORAGE_DRIVER": "ocis",
            "UPLOAD_DELETE_WAIT_TIME": "0",
            "EMAIL_HOST": "localhost",
            "EMAIL_PORT": EMAIL_PORT,
        }
        behat_env.update(cfg["extraEnvironment"])

        print(f"Running suites: {behat_suites_raw} (type: {acceptance_test_type})")
        result = subprocess.run(
            ["make", "-C", str(repo_root), "test-acceptance-api"],
            env=behat_env,
        )
        return result.returncode

    finally:
        cleanup()


if __name__ == "__main__":
    sys.exit(main())