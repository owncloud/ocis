# Documentation
The documentation is incomplete because the Collaboration server is in a development state.

# Infinite Scale Collaboration Deployment Example

This deployment example of the oCIS with the new Collaboration server.

## Overview

* oCIS, Collaboration server, Collabora or OnlyOffice running behind Traefik as reverse proxy
* Collabora or OnlyOffice enable you to edit documents in your browser
* Collaboration server acts as a bridge to make the oCIS storage accessible to Collabora or OnlyOffice
* Traefik generating self-signed certificates for local setup or obtaining valid SSL certificates for a server setup
Please note: Against the stack that uses [wopiserver](https://owncloud.dev/ocis/deployment/ocis_wopi/), we don't need the app_provider anymore. The new Collaboration server now includes an app_provider.

### Running

```bash
docker compose -f docker-compose.collabora.yml up -d
```

```bash
docker compose -f docker-compose.onlyoffice.yml up -d
```

Also see the [Admin Documentation](https://doc.owncloud.com/ocis/latest/index.html) for administrative and more configuration details.
