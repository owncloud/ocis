Security: Bump golang.org/x/text to v0.39.0

We bumped golang.org/x/text to v0.39.0 to address GO-2026-5970, an infinite loop
on invalid input. This unblocks the govulncheck CI check, which fails on fixable
vulnerabilities.
