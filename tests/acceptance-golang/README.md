# Pre-requisites
 To run this project you need to have the following installed:
 - [Go](https://go.dev/doc/install) (1.18 or higher)
 - [godog](https://github.com/cucumber/godog)
 ```bash
 go install github.com/cucumber/godog/cmd/godog@latest
```
- Add the path to the godog binary to your PATH environment variable
```bash
export PATH=$PATH:$GOPATH/bin
```
- Run `OCIS` server

# Run acceptance tests
- Change directory to `tests/acceptance-golang`
- Install all the test dependencies with `go mod vendor` command
- Change value of `baseUrl` in `test-helpers/setup-helper.go` to the URL of the OCIS server
- Run the tests with `godog` command
```feature
godog run <path-to-feature-file>

godog run features/create-user.feature
```
