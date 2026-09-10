Change: Replace axios with the native fetch API in Web

The Web frontend no longer depends on axios. All HTTP traffic goes through a
single fetch-based client, and the libre-graph client is generated from the
`typescript-fetch` template instead of `typescript-axios`.

Sending requests is a drop-in change. `HttpClient` keeps its per-request config
including `data`, keeps `cancel()`, still resolves a non-JSON body as text, and
still exposes response headers as `headers['etag']` as well as
`headers.get('etag')`. Errors still carry the body and status as `error.data` /
`error.statusCode` and as `error.response.data` / `error.response.status`.
`mockAxiosResolve` and `mockAxiosReject` still work, deprecated in favour of
`mockHttpResponse` and `mockHttpError`.

Per-request `headers` are merged over the client-wide ones case-insensitively, so
an override replaces the header it names whatever case either side used.

Code that touches axios directly has to be adapted:

- `new HttpClient()` takes `{ baseUrl, staticHeaders, headers, onResponse }`
  instead of `{ config, requestInterceptor, responseInterceptor }`. Clients from
  `ClientService` are unaffected.
- The per-request config drops `timeout`, `withCredentials`, `onUploadProgress`,
  `cancelToken`, `paramsSerializer`, `transformRequest` / `transformResponse`,
  `validateStatus` and `baseURL`; `responseType` drops `document` and `stream`.
- The per-request `headers` are typed as `HeadersInit`, so a `Headers` is accepted
  as well as a plain object. Assigning into them after the fact
  (`config.headers.Authorization = …`) no longer type-checks; build the object
  first, or use `new Headers()` and `set()`.
- `error.response` carries the body on `error.response.data`, as it did with
  axios. Its underlying stream is consumed, so `error.response.json()` is not
  available.
- `graph()`, `ocs()`, `UrlSign` and `WebDavOptions` take a `FetchClient`, and
  the latter two rename `axiosClient` to `httpClient`. `webdav()` is unchanged.
- Graph fields with an OData annotation use their generated camelCase names,
  e.g. `atLibreGraphPermissionsActions`. The wire format is unchanged.
- The generated client loses its `*ApiFactory`, `*ApiFp` and
  `*AxiosParamCreator` exports. The `*Api` classes now take one options object
  per operation and resolve with the payload.
- `@ownclouders/web-test-helpers` no longer has a `mocks/axios` module path.

https://github.com/owncloud/ocis/pull/12910
