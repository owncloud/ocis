# sbomgen

Generates the Software Bill of Materials (SPDX 3.0.1, JSON-LD) that is attached to every oCIS release.

It scans the Go modules (`go mod download`, `go mod graph`) and the pnpm lockfile of `web`
(production dependencies only, unless `-npm-dev` is set), reads the license and copyright statements of
every package, and converts the result with the official
[SPDX tools-golang](https://github.com/spdx/tools-golang) library. The document is validated before it is
written; an invalid document makes the tool exit with an error.

```
go -C tools/sbomgen build -o /tmp/sbomgen .
/tmp/sbomgen \
  -go-module . \
  -pnpm-lock web/pnpm-lock.yaml \
  -reuse-toml REUSE.toml \
  -version 9.0.0 \
  -out ocis-9.0.0.spdx3.json
```

Run it from the repository root. Needs network access (Go module proxy, npm registry, spdx.org license list).
Use `-npm-cache <file>` to keep the npm license facts between runs; `-help` lists all flags.

The file name must end in `.spdx3.json`, otherwise SPDX tools such as the online validator read it as SPDX 2.x.

Trivy produces the CycloneDX file that is released next to it (see `release.yml`).

Not covered: binaries and container images (the SBOM lists source dependencies), and a signature of the SBOM.
