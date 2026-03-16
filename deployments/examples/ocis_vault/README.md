---
document this deployment example in: docs/ocis/deployment/ocis_vault.md
---

Please refer to [our documentation](https://owncloud.dev/ocis/deployment/ocis_vault/)
for instructions on how to deploy this scenario.

## Local web development

The `ocis` service mounts the web repo from the host (`/Users/mk/dev/kiteworks/web`)
and serves its `dist/` via `WEB_ASSET_CORE_PATH`. oCIS generates `config.json`
dynamically from its own config — the one in `dist/` is ignored.

Run the web build in watch mode:

```bash
cd /Users/mk/dev/kiteworks/web
pnpm build:w
```

This is `vite build --watch` — a full production rebuild on every save (no HMR,
takes 20-60 s). Hard-refresh the browser after each rebuild to pick up changes.
No container restart needed.
