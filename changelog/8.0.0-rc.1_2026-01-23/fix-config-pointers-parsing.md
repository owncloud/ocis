Bugfix: Support pointer types in config environment variable decoding

Added support for decoding pointer types (*bool, *int, *string, etc.) in the envdecode package, allowing configuration fields to distinguish between unset (nil) and explicitly set values. Changed `WEB_OPTION_EMBED_ENABLED` from string to *bool type to enable explicit false values.

https://github.com/owncloud/ocis/pull/11815
