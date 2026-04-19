# ownCloud Infinite Scale

<!-- OSPO-managed README | Generated: 2026-04-16 | v2 -->

[![License](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](LICENSE) [![ownCloud OSPO](https://img.shields.io/badge/OSPO-ownCloud-blue)](https://kiteworks.com/opensource) [![Docker Hub](https://img.shields.io/docker/pulls/owncloud)](https://hub.docker.com/r/owncloud/ocis)

ownCloud Infinite Scale (oCIS) is the next-generation file sync, share, and collaboration platform built in Go. It is a single-binary, cloud-native server that replaces the classic PHP-based ownCloud Server with a microservices architecture supporting S3-compatible storage backends, spaces-based file organization, OpenID Connect authentication, and the Libre Graph API for extensibility.

## Getting Started

Follow the steps below to deploy oCIS using Docker or build from source.

### Docker Quickstart

```bash
mkdir -p $HOME/ocis/ocis-config $HOME/ocis/ocis-data
docker run --rm -it \
    --mount type=bind,source=$HOME/ocis/ocis-config,target=/etc/ocis \
    --mount type=bind,source=$HOME/ocis/ocis-data,target=/var/lib/ocis \
    -p 9200:9200 \
    owncloud/ocis init
docker run -d \
    --mount type=bind,source=$HOME/ocis/ocis-config,target=/etc/ocis \
    --mount type=bind,source=$HOME/ocis/ocis-data,target=/var/lib/ocis \
    -p 9200:9200 \
    owncloud/ocis server
```

### Build from Source

```bash
make -C ocis build
./ocis/bin/ocis init
./ocis/bin/ocis server
```

### Run Tests

```bash
make test
```

## Documentation

- [Admin Documentation](https://doc.owncloud.com/ocis/next/)
- [Developer Documentation](https://owncloud.dev/)
- [Deployment Guide](https://doc.owncloud.com/ocis/next/deployment/)
- [oCIS Quickstart with Docker](https://owncloud.dev/ocis/guides/ocis-mini-eval/)

## Features

Key capabilities of ownCloud Infinite Scale:

### Clients

oCIS supports all official ownCloud clients:

- [ownCloud Web](https://github.com/owncloud/web)
- [Android](https://github.com/owncloud/android)
- [iOS](https://github.com/owncloud/ios-app)
- [Desktop](https://github.com/owncloud/client/)

### Web Office Integration

Collaborative editing via [WOPI](https://github.com/cs3org/wopiserver) with:

- [Collabora Online](https://github.com/CollaboraOnline/online)
- [OnlyOffice Docs](https://github.com/ONLYOFFICE/DocumentServer)
- [Microsoft Office Online Server](https://owncloud.com/microsoft-office-online-integration-with-wopi/)

### Authentication

Users authenticate via [OpenID Connect](https://openid.net/connect/) using either an external IdP (e.g., [Keycloak](https://www.keycloak.org/)) or the embedded [LibreGraph Connect](https://github.com/libregraph/lico) identity provider.

### Architecture

oCIS is delivered as a single binary or container with a microservices architecture built on [reva](https://reva.link/). It scales from a Raspberry Pi to a Kubernetes cluster and uses open APIs including [WebDAV](http://www.webdav.org/) and [CS3](https://github.com/cs3org/cs3apis/). No external database or IdP is required for basic deployments.

### Building from Source

Requires Go >= 1.25.7 and a C compiler (for reva's C-Go dependencies):

```bash
git clone git@github.com:owncloud/ocis.git && cd ocis
make generate
make -C ocis build
./ocis/bin/ocis init
IDM_CREATE_DEMO_USERS=true ./ocis/bin/ocis server
```

Access the web UI at `http://localhost:9200`.

## Part of ownCloud Infinite Scale

This is the core repository for oCIS -- the primary product of the ownCloud open source project.

- [Download latest release](https://download.owncloud.com/ocis/ocis/stable/?sort=time&order=desc)
- [Docker Hub: owncloud/ocis](https://hub.docker.com/r/owncloud/ocis)

## Community & Support

**[Star](https://github.com/owncloud/ocis)** this repo and **Watch** for release notifications!

- [ownCloud Website](https://owncloud.com)
- [Community Discussions](https://github.com/orgs/owncloud/discussions)
- [Matrix Chat](https://app.element.io/#/room/#owncloud:matrix.org)
- [Documentation](https://doc.owncloud.com)
- [Enterprise Support](https://owncloud.com/contact-us/)
- [OSPO Home](https://kiteworks.com/opensource)

## Contributing

We welcome contributions! Please read the [Contributing Guidelines](CONTRIBUTING.md)
and our [Code of Conduct](CODE_OF_CONDUCT.md) before getting started.

### Workflow

- **Rebase Early, Rebase Often!** We use a rebase workflow. Always rebase on the target branch before submitting a PR.
- **Dependabot**: Automated dependency updates are managed via Dependabot. Review and merge dependency PRs promptly.
- **Signed Commits**: All commits **must** be PGP/GPG signed. See [GitHub's signing guide](https://docs.github.com/en/authentication/managing-commit-signature-verification).
- **DCO Sign-off**: Every commit must carry a `Signed-off-by` line:
  ```
  git commit -s -S -m "your commit message"
  ```
- **GitHub Actions Policy**: Workflows may only use actions that are (a) owned by `owncloud`, (b) created by GitHub (`actions/*`), or (c) verified in the GitHub Marketplace.

## Translations

Help translate this project on Transifex:
**<https://explore.transifex.com/owncloud-org/owncloud-web/>**

Please submit translations via Transifex -- do not open pull requests for translation changes.

## Security

**Do not open a public GitHub issue for security vulnerabilities.**

Report vulnerabilities at **<https://security.owncloud.com>** -- see [SECURITY.md](SECURITY.md).

Bug bounty: [YesWeHack ownCloud Program](https://yeswehack.com/programs/owncloud-bug-bounty-program)

## License

This project is licensed under the [Apache-2.0](LICENSE).

## About the ownCloud OSPO

The [Kiteworks Open Source Program Office](https://kiteworks.com/opensource), operating under
the [ownCloud](https://owncloud.com) brand, launched on May 5, 2026, to steward the open source
ecosystem around ownCloud's products. The OSPO ensures transparent governance, license compliance,
community health, and sustainable collaboration between the open source community and
[Kiteworks](https://www.kiteworks.com), which acquired ownCloud in 2023.

- **OSPO Home**: <https://kiteworks.com/opensource>
- **GitHub**: <https://github.com/owncloud>
- **ownCloud**: <https://owncloud.com>

For questions about the OSPO or licensing, contact ospo@kiteworks.com.

> **License status:** This repository is already licensed under Apache-2.0 -- the OSPO target license.
> No migration is required.
