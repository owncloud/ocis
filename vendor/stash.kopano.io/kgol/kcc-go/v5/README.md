# Kopano Core client for Go

This implements a minimal client interfacing with a couple of
SOAP methods of a Kopano server.

## Quickstart

Make sure you have Go 1.13 or later installed. This project uses Go modules.

```
mkdir -p ~/go/src/stash.kopano.io/kgol
cd ~/go/src/stash.kopano.io/kgol
git clone <THIS-PROJECT> kcc-go
cd kcc-go
```

All the rest of this document assumes you have installed like above and pwd is
the `kcc-go` directory.

## Environment variables

| Environment variable       | Description                                   |
|----------------------------|-----------------------------------------------|
| KOPANO_SERVER_DEFAULT_URI  | URI used to connect to Kopano server          |
| TEST_USERNAME              | Kopano username used in unit tests            |
| TEST_PASSWORD              | Kopano username's password used in unit tests |

## Testing

Running the unit tests requires a Kopano Server with accessible SOAP service.
Make sure to set the environment variables as listed above to match your Kopano
server details.

Testing requires a running Kopano Groupware Storage server.

```
go test -v
```

## Benchmark

For testing there is also a benchmark test.

```
go test -v -bench=. -run BenchmarkLogon -benchmem
BenchmarkLogon-8            2000            591907 ns/op           20509 B/op       217 allocs/op
PASS
ok      stash.kopano.io/kc/kcc-go       1.255s
```

## Test server

For example usage, a simple test HTTP server `kuserd` is included. Run it like
this:

```
glide install
go install -v ./cmd/kuserd && KOPANO_USERNAME=system KOPANO_PASSWORD= kuserd serve
INFO[0000] serve started
```

Make sure to specify `KOPANO_USERNAME` and `KOPANO_PASSWORD` according to your
setup. It must be a valid existing user. If not give, the server defaults to the
`SYSTEM` user with empty password.

### Endpoints

The `kuserd` test server exposes a bunch of endpoints for easy testing with
curl or similar.

#### /logon
```
curl -v "http://user1:pass@127.0.0.1:8769/logon"
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 8769 (#0)
* Server auth using Basic with user 'user1'
> GET /logon HTTP/1.1
> Host: 127.0.0.1:8769
> Authorization: Basic dXNlcjE6cGFzcw==
> User-Agent: curl/7.47.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Tue, 07 Nov 2017 18:16:07 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 127.0.0.1 left intact
```

#### /logoff?id=${sessionID}
```
curl -v "http://127.0.0.1:8769/logoff?id=17227711350321618063"
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 8769 (#0)
> GET /logoff?id=17227711350321618063 HTTP/1.1
> Host: 127.0.0.1:8769
> User-Agent: curl/7.47.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Tue, 07 Nov 2017 19:16:25 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 127.0.0.1 left intact
```

#### /userinfo?username=${username}

```
curl -v "http://127.0.0.1:8769/userinfo?username=system"
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 8769 (#0)
> GET /userinfo?username=system HTTP/1.1
> Host: 127.0.0.1:8769
> User-Agent: curl/7.47.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 07 Nov 2017 18:16:28 GMT
< Content-Length: 226
<
{
  "ulUserID": 2,
  "lpszUsername": "SYSTEM",
  "lpszMailAddress": "postmaster@localhost",
  "lpszFullName": "SYSTEM",
  "ulIsAdmin": 2,
  "ulIsNonActive": 0,
  "sUserId": "AAAAAKwhqVBA0+5Isxn7p1MwRCUAAAAABgAAAAIAAAAAAAAA"
}
* Connection #0 to host 127.0.0.1 left intact
```

#### /error?er=${error_code}

Converts Kopano Core error codes to a meaningful string.

```
curl "http://127.0.0.1:8769/error?er=0x80000002"
KCERR_NOT_FOUND: (Not Found) (KC:0x80000002)
```

#### /errors

Lists all known Errors with integer and hex representation codes.

### Benchmark / load tests

Use [hey](https://github.com/rakyll/hey) to test it.

```
hey -n 10000 -c 200 -a user1:pass http://127.0.0.1:8769/logon
Summary:
  Total:        1.9587 secs
  Slowest:      0.3196 secs
  Fastest:      0.0006 secs
  Average:      0.0376 secs
  Requests/sec: 5105.3052

Response time histogram:
  0.001 [1]     |
  0.032 [2035]  |∎∎∎∎∎∎∎∎∎∎
  0.064 [7774]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.096 [120]   |∎
  0.128 [50]    |
  0.160 [16]    |
  0.192 [1]     |
  0.224 [0]     |
  0.256 [1]     |
  0.288 [0]     |
  0.320 [2]     |

Latency distribution:
  10% in 0.0289 secs
  25% in 0.0333 secs
  50% in 0.0364 secs
  75% in 0.0402 secs
  90% in 0.0460 secs
  95% in 0.0537 secs
  99% in 0.0850 secs

Details (average, fastest, slowest):
  DNS+dialup:    0.0003 secs, 0.0000 secs, 0.0880 secs
  DNS-lookup:    0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:     0.0000 secs, 0.0000 secs, 0.0936 secs
  resp wait:     0.0367 secs, 0.0005 secs, 0.2742 secs
  resp read:     0.0004 secs, 0.0000 secs, 0.0458 secs

Status code distribution:
  [200] 10000 responses
```

```
hey -n 10000 -c 200 "http://127.0.0.1:8769/userinfo?username=system"
Summary:
  Total:        2.3543 secs
  Slowest:      1.0623 secs
  Fastest:      0.0004 secs
  Average:      0.0445 secs
  Requests/sec: 4247.5039
  Total data:   2260000 bytes
  Size/request: 226 bytes

Response time histogram:
  0.000 [1]     |
  0.107 [9928]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.213 [46]    |
  0.319 [23]    |
  0.425 [0]     |
  0.531 [0]     |
  0.638 [0]     |
  0.744 [0]     |
  0.850 [0]     |
  0.956 [0]     |
  1.062 [2]     |

Latency distribution:
  10% in 0.0265 secs
  25% in 0.0339 secs
  50% in 0.0418 secs
  75% in 0.0522 secs
  90% in 0.0649 secs
  95% in 0.0737 secs
  99% in 0.0965 secs

Details (average, fastest, slowest):
  DNS+dialup:    0.0003 secs, 0.0000 secs, 0.0872 secs
  DNS-lookup:    0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:     0.0000 secs, 0.0000 secs, 0.0270 secs
  resp wait:     0.0441 secs, 0.0004 secs, 1.0623 secs
  resp read:     0.0001 secs, 0.0000 secs, 0.0120 secs

Status code distribution:
  [200] 10000 responses
```
