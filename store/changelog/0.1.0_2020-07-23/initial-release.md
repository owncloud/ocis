Change: Initial release of basic version

We have built a new service which implements go micro's [store interface](https://github.com/micro/development/blob/master/design/framework/store.md),
so that we can provide a key-value-store through GRPC for our services. The underlying implementation stores data as json
files on disk and maintains an index using [bleve](https://github.com/blevesearch/bleve).

<https://github.com/owncloud/ocis/store/pull/1>
<https://github.com/owncloud/ocis/store/pull/2>
