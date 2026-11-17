Bugfix: Filter spurious TLS handshake errors from CI test runner output

Go's `net/http` server emits `TLS handshake error from <addr>: ...` to stderr
whenever a client connects to an HTTPS port but does not complete the handshake.
This happens routinely in acceptance test runs: health-check probes, Playwright
browser initialisation, and curl readiness polls all trigger the message even
though no test ever fails because of it.  The messages accumulate in the
hundreds per run and push genuinely interesting output far down the CI log.

The five Python test-runner scripts (`run-github.py`, `run-litmus.py`,
`run-cs3api.py`, `run-wopi.py`) now pipe each long-lived ocis/ociswrapper
subprocess's stderr through a background thread that drops lines containing
`TLS handshake error` and forwards everything else to the runner's stderr
unchanged.  Real errors remain fully visible.

The two Bash scripts (`tests/ldap-smoke/run.sh`,
`tests/ldap-smoke-distributed/run.sh`) use a Bash process substitution
(`2> >(grep -v "TLS handshake error" >&2)`) for the same effect without
requiring a log file.

https://github.com/owncloud/ocis/pull/13499
