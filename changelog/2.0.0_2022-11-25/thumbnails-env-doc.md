Enhancement: Add description tags to the thumbnails config structs

Added description tags to the config structs in the thumbnails service so they will be included in the config documentation.

**Important**
If you ran `ocis init` with the `v2.0.0-alpha*` version then you have to manually add the `transfer_secret` to the ocis.yaml.

Just open the `ocis.yaml` config file and look for the thumbnails section.
Then add a random `transfer_secret` so that it looks like this:

```yaml
thumbnails:
  thumbnail:
    transfer_secret: <put random value here>
```

https://github.com/owncloud/ocis/pull/3752
