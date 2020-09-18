Enhancement: update reva to v0.1.1-0.20200701152626-2f6cc60e2f66

- Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66 (#328)
- Use sync.Map on pool package (reva/#909)
- Use mutex instead of sync.Map (reva/#915)
- Use gatewayProviders instead of storageProviders on conn pool (reva/#916)
- Add logic to ls and stat to process arbitrary metadata keys (reva/#905)
- Preliminary implementation of Set/UnsetArbitraryMetadata (reva/#912)
- Make datagateway forward headers (reva/#913, reva/#926)
- Add option to cmd upload to disable tus (reva/#911)
- OCS Share Allow date-only expiration for public shares (#288, reva/#918)
- OCS Share Remove array from OCS Share update response (#252, reva/#919)
- OCS Share Implement GET request for single shares  (#249, reva/#921)

https://github.com/owncloud/ocis/ocis-revapull/328
https://github.com/cs3org/reva/pull/909
https://github.com/cs3org/reva/pull/915
https://github.com/cs3org/reva/pull/916
https://github.com/cs3org/reva/pull/905
https://github.com/cs3org/reva/pull/912
https://github.com/cs3org/reva/pull/913
https://github.com/cs3org/reva/pull/926
https://github.com/cs3org/reva/pull/911
https://github.com/owncloud/ocis/ocis-revaissues/288
https://github.com/cs3org/reva/pull/918
https://github.com/owncloud/ocis/ocis-revaissues/252
https://github.com/cs3org/reva/pull/919
https://github.com/owncloud/ocis/ocis-revaissues/249
https://github.com/cs3org/reva/pull/921

