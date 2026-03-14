Bugfix: Fix IDP build on FreeBSD by disabling absolute Babel runtime

The `babel-preset-react-app` preset defaults to `absoluteRuntime: true`,
which hardcodes absolute paths to `@babel/runtime` helpers. These paths
fail to resolve on non-Linux platforms like FreeBSD. Setting
`absoluteRuntime: false` makes Babel resolve the runtime relative to
the source file, which works across all platforms.

https://github.com/owncloud/ocis/pull/XXXX
https://github.com/owncloud/ocis/issues/12065
