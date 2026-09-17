Enhancement: Bump dependencies

Bumped Go and npm dependencies, including security fixes:

- `github.com/owncloud/reva/v2` to `v2.0.0-20260519092700-9da01c6fb954`
- `github.com/shamaton/msgpack/v2` v2.4.0 → v2.4.1 (CVE: denial of service)
- `filippo.io/edwards25519` v1.1.0 → v1.1.1
- `github.com/cloudflare/circl` v1.6.1 → v1.6.3
- `github.com/russellhaering/goxmldsig` v1.5.0 → v1.6.0
- `postcss`, `fast-uri`, `@babel/plugin-transform-modules-systemjs` (npm, via pnpm lock regen)
- GitHub Actions: `actions/upload-artifact` 4→7, `actions/download-artifact` 4→8, `pnpm/action-setup` 5→6, `fpicalausa/remove-stale-branches` 2.4→2.6

https://github.com/owncloud/ocis/pull/12325
