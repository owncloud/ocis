Security: Upgrade Go to 1.25.12

Bumped the Go toolchain used to build the release binaries and Docker images
from 1.25.11 to 1.25.12. Go 1.25.11 is affected by CVE-2026-39822 (os.Root
symlink following allows directory traversal), which is fixed in 1.25.12 and
was blocking the release image security scan.

https://github.com/owncloud/ocis/pull/12602
