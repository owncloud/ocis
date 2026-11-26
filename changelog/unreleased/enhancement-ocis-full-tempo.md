Enhancement: Add Grafana Tempo to the ocis_full example

The `ocis_full` example now has an optional setup for tracing. It starts
Grafana Tempo and Grafana. oCIS sends the traces to Tempo. Grafana shows them.
To enable the setup, set `TEMPO=:monitoring_tracing/tempo.yml` in the `.env`
file. Use the setup only for tests and debugging.

https://github.com/owncloud/ocis/pull/13605
