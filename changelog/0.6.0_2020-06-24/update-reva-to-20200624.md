Enhancement: update reva to v0.1.1-0.20200624063447-db5e6635d5f0

- Updated reva to v0.1.1-0.20200624063447-db5e6635d5f0 (#279)
- Local storage: URL-encode file ids to ease integration with other microservices like WOPI (reva/#799)
- Mentix fixes (reva/#803, reva/#817)
- OCDAV: fix returned timestamp format (#116, reva/#805)
- OCM: add default prefix (#814)
- add the content-length header to the responses (reva/#816)
- Deps: clean (reva/#818)
- Fix trashbin listing (#112, #253, #254, reva/#819)
- Make the json publicshare driver configurable (reva/#820)
- TUS: Return metadata headers after direct upload (ocis/#216, reva/#813)
- Set mtime to storage after simple upload (#174, reva/#823, reva/#841)
- Configure grpc client to allow for insecure conns and skip server certificate verification (reva/#825)
- Deployment: simplify config with more default values (reva/#826, reva/#837, reva/#843, reva/#848, reva/#842)
- Separate local fs into home and with home disabled (reva/#829)
- Register reflection after other services (reva/#831)
- Refactor EOS fs (reva/#830)
- Add ocs-share-permissions to the propfind response (#47, reva/#836)
- OCS: Properly read permissions when creating public link (reva/#852)
- localfs: make normalize return associated error (reva/#850)
- EOS grpc driver (reva/#664)
- OCS: Add support for legacy public link arg publicUpload (reva/#853)
- Add cache layer to user REST package (reva/#849)
- Meshdirectory: pass query params to selected provider (reva/#863)
- Pass etag in quotes from the fs layer (#269, reva/#866, reva/#867)
- OCM: use refactored cs3apis provider definition (reva/#864)

https://github.com/owncloud/ocis-reva/pull/279
https://github.com/owncloud/cs3org/reva/pull/799
https://github.com/owncloud/cs3org/reva/pull/803
https://github.com/owncloud/cs3org/reva/pull/817
https://github.com/owncloud/ocis-reva/issues/116
https://github.com/owncloud/cs3org/reva/pull/805
https://github.com/owncloud/cs3org/reva/pull/814
https://github.com/owncloud/cs3org/reva/pull/816
https://github.com/owncloud/cs3org/reva/pull/818
https://github.com/owncloud/ocis-reva/issues/112
https://github.com/owncloud/ocis-reva/issues/253
https://github.com/owncloud/ocis-reva/issues/254
https://github.com/owncloud/cs3org/reva/pull/819
https://github.com/owncloud/cs3org/reva/pull/820
https://github.com/owncloud/ocis/issues/216
https://github.com/owncloud/ocis-reva/issues/174
https://github.com/owncloud/cs3org/reva/pull/823
https://github.com/owncloud/cs3org/reva/pull/841
https://github.com/owncloud/cs3org/reva/pull/813
https://github.com/owncloud/cs3org/reva/pull/825
https://github.com/owncloud/cs3org/reva/pull/826
https://github.com/owncloud/cs3org/reva/pull/837
https://github.com/owncloud/cs3org/reva/pull/843
https://github.com/owncloud/cs3org/reva/pull/848
https://github.com/owncloud/cs3org/reva/pull/842
https://github.com/owncloud/cs3org/reva/pull/829
https://github.com/owncloud/cs3org/reva/pull/831
https://github.com/owncloud/cs3org/reva/pull/830
https://github.com/owncloud/ocis-reva/issues/47
https://github.com/owncloud/cs3org/reva/pull/836
https://github.com/owncloud/cs3org/reva/pull/852
https://github.com/owncloud/cs3org/reva/pull/850
https://github.com/owncloud/cs3org/reva/pull/664
https://github.com/owncloud/cs3org/reva/pull/853
https://github.com/owncloud/cs3org/reva/pull/849
https://github.com/owncloud/cs3org/reva/pull/863
https://github.com/owncloud/ocis-reva/issues/269
https://github.com/owncloud/cs3org/reva/pull/866
https://github.com/owncloud/cs3org/reva/pull/867
https://github.com/owncloud/cs3org/reva/pull/864
