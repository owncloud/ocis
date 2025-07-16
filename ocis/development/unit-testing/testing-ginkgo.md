---
title: "Testing with Ginkgo"
date: 2024-04-25T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development/unit-testing
geekdocFilePath: testing-ginkgo.md

---

{{< toc >}}

In this section we try to enable developers to write tests in oCIS using Ginkgo and Gomega and explain how to mock other microservices to also cover some integration tests. The full documentation of the tools can be found on the [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/) websites.

{{% hint type=tip icon=gdoc_link title="Reading the documentation" %}}
This page provides only a basic introduction to get started with Ginkgo and Gomega. For more detailed information, please refer to the official documentation.

**Useful Links:**

- [Ginkgo](https://onsi.github.io/ginkgo/)
- [Gomega](https://onsi.github.io/gomega/)
- [Mockery](https://vektra.github.io/mockery/latest/)

{{% /hint %}}

## Prerequisites

To use Ginkgo, you need to install the Ginkgo CLI. You can install it using the following command:

```bash
go install github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega/...
```

## Getting Started

Navigate to the directory where you want to write your tests and run the following command:

### Bootstrap

```bash
cd ocis/ocis-pkg/config/parser
ginkgo bootstrap
Generating ginkgo test suite bootstrap for parser in:
        parser_suite_test.go

```

This command creates a `parser_suite_test.go` file in the parser directory. This file contains the test suite for the parser package.

```go
package parser_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestParser(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Parser Suite")
}
```

Ginkgo defaults to setting up the suite as a `*_test` package to encourage you to only test the external behavior of your package, not its internal implementation details.

After the package `parser_test` declaration we import the ginkgo and gomega packages into the test's top-level namespace by performing a `.` dot-import. Since Ginkgo and Gomega are DSLs (Domain-specific Languages) this makes the tests more natural to read. If you prefer, you can avoid the dot-import via `ginkgo bootstrap --nodot`. Throughout this documentation we'll assume dot-imports.

With the bootstrap complete, you can now run your tests using the `ginkgo` command:

```bash
ginkgo

Running Suite: Parser Suite - <local-path>/ocis/ocis-pkg/config/parser
===============================================================================================
Random Seed: 1714076559

Will run 0 of 0 specs

Ran 0 of 0 Specs in 0.000 seconds
SUCCESS! -- 0 Passed | 0 Failed | 0 Pending | 0 Skipped
PASS

Ginkgo ran 1 suite in 7.0058606s
Test Suite Passed
```

Under the hood, ginkgo is simply calling `go test`. While you can run `go test` instead of the ginkgo CLI, Ginkgo has several capabilities that can only be accessed via `ginkgo`. We generally recommend users embrace the ginkgo CLI and treat it as a first-class member of their testing toolchain.

### Adding Specs to the Suite

```bash
ginkgo generate parser
Generating ginkgo test for Parser in:                                                                                                                                                   ✔  7s  22:22:46 
  parser_test.go
```

This will generate a `parser_test.go` file in the parser directory. This file contains the test suite for the parser package.

```go
package parser_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
)

var _ = Describe("Parser", func() {

})
```

## Writing Specs

### Describe

The `Describe` block is used to describe the behavior of a particular component of your code. It is a way to group together related specs. The `Describe` block takes a string and a function. The string is a description of the component you are describing, and the function contains the specs that describe the behavior of that component.

```go
var _ = Describe("Parser", func() {
    // Specs go here
})
```

### Context

The `Context` block is used to further describe the behavior of a component. It is a way to group together related specs within a `Describe` block. The `Context` block takes a string and a function. The string is a description of the context you are describing, and the function contains the specs that describe the behavior of that context.

```go
var _ = Describe("Parser", func() {
    Context("when the input is valid", func() {
        // Specs go here
    })
})
```

### It

The `It` block is used to describe a single spec. It takes a string and a function. The string is a description of the behavior you are specifying, and the function contains the code that exercises that behavior.

```go
var _ = Describe("Parser", func() {
    Context("when the input is valid", func() {
        It("parses the input", func() {
            // Spec code goes here
        })
    })
})
```

### Expect

The `Expect` function is used to make assertions in your specs. It takes a value and returns an `*Expectation`. You can then chain methods on the `*Expectation` to make assertions about the value.

```go
var _ = Describe("Parser", func() {
    Context("when the input is valid", func() {
        It("parses the input", func() {
            result := parser.Parse("valid input")
            Expect(result).To(Equal("expected output"))
        })
    })
})
```

### BeforeEach

The `BeforeEach` block is used to run a setup function before each spec in a `Describe` or `Context` block. It takes a function that contains the setup code.

```go
package parser_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/owncloud/ocis/v2/ocis-pkg/config"

    p "github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
)

var _ = Describe("Parser", func() {
    var c *config.Config

    BeforeEach(func() {
        c = config.DefaultConfig()
    })

    Context("when the input is valid", func() {
        It("parses the input", func() {
            err := p.ParseConfig(c, false)
            Expect(err).ToNot(HaveOccurred())
            Expect(c.Commons.OcisURL).To(Equal("https://localhost:9200"))
        })
    })
})
```

Let us take a closer look at the code above:

We are following the recommended practise on variables to **"declare in container nodes"** and **"initialize in setup nodes"**. This is why we are declaring the `c` variable at the top of the `Describe` block and initializing it in the `BeforeEach` block. This is important to get isolated test steps which can be run in any order and even in parallel.

Let us take a look at a bad example where we are polluting the spec by not following this recommended practise:

```go
package parser_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/owncloud/ocis/v2/ocis-pkg/config"

    p "github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
)


var _ = Describe("Parser", func() {
    c := config.DefaultConfig()

    Context("when the defaults are applied", func() {
        It("fails to parse the input", func() {
            c.TokenManager.JWTSecret = "" // bam! we have changed the closure variable and it will never be reset
            err := p.ParseConfig(c, false)
            Expect(err).To(HaveOccurred())
        })
        It("parses the input", func() {
            err := p.ParseConfig(c, false)
            Expect(err).ToNot(HaveOccurred())
            Expect(c.Commons.OcisURL).To(Equal("https://localhost:9200"))
        })
    })
})
```

{{% hint type="warning" title="Specs MUST be clean and independent"%}}
Always **declare variables in the container node**(which are basically `Describe()` and `Context()`)

and **initialize your variables in the setup nodes.** (which are basically `BeforeEach()` and `JustBeforeEach()`).

This will ensure that your specs are clean and independent of each other.
{{% /hint %}}

### Focused Specs

You can focus on a single spec by adding an `F` in front of the `It` block. This will run only the focused spec.

```go
var _ = Describe("Parser", func() {
    Context("when the input is valid", func() {
        FIt("parses the input", func() {
            result := parser.Parse("valid input")
            Expect(result).To(Equal("expected output"))
        })
    })
})
```

### Pending Specs

You can mark a spec as pending by adding a `P` in front of the `It` block. This will skip the spec.

```go
var _ = Describe("Parser", func() {
    Context("when the input is valid", func() {
        PIt("parses the input", func() {
            result := parser.Parse("valid input")
            Expect(result).To(Equal("expected output"))
        })
    })
})
```

### Test Driven Development

You can run the tests in watch mode to follow a test-driven development approach. This will run the tests every time you save a file.

```bash
ginkgo watch
```

## Mocking

In oCIS, we use the `mockery` tool to generate mocks for interfaces. [Mockery](https://vektra.github.io/mockery/latest/) is a simple tool that generates mock implementations of Go interfaces. It is useful for writing tests against interfaces instead of concrete types. We can use it to mock requests to other microservices to cover some integration tests. We should already have a number of mocks in the project. The mocks are configured on the packages level in the `.mockery.yaml` files.

**Example file:**

```yaml
with-expecter: true
filename: "{{.InterfaceName | snakecase }}.go"
dir: "{{.PackageName}}/mocks"
mockname: "{{.InterfaceName}}"
outpkg: "mocks"
packages:
    github.com/owncloud/ocis/v2/ocis-pkg/oidc:
        interfaces:
            OIDCClient:
```

We should add missing mocks to this file and define the interfaces we want to mock. After that, we can generate the mocks by running `mockery` in the repo, it will find all the `.mockery.yaml` files and generate the mocks for the interfaces defined in them.

Our mocks are generated with the setting `with-expecter: true`. This allows us to use type-safe methods to generate the call expectations by simply calling `EXPECT()` on the mock object.

{{% hint type="tip" title="Type safe mock identifiers" %}}
By using `EXPECT()` on the mock object, we can work with type-safe methods to generate the call expectations.
{{% /hint %}}

**Example of a mocked gateway client**

In our oCIS services we need to use a gateway pool selector to get the gateway client.

We should always use the constructor on a new mock like `gatewayClient = cs3mocks.NewGatewayAPIClient(GinkgoT())`. This brings us two advantages:

- The `AssertExpectations` method is registered to be called at the end of the tests via `t.Cleanup()` method.
- The `testing.TB` interface is registered on the `mock.Mock` so that tests don't panic when a call on the mock is unexpected.

```go
package publicshareprovider_test

import (
    "context"
    "time"


    "github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
    cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "google.golang.org/grpc"
)

var _ = Describe("PublicShareProvider", func() {
    // declare in container nodes
    var (
        gatewayClient   *cs3mocks.GatewayAPIClient
        gatewaySelector pool.Selector
    )

    BeforeEach(func() {
        // initialize in setup nodes
        pool.RemoveSelector("GatewaySelector" + "any")
        // create a new mock client
        gatewayClient = cs3mocks.NewGatewayAPIClient(GinkgoT())
        gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
            "GatewaySelector",
            "any",
            func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
                return gatewayClient
            },
        )
    })
    Context("The user has the permission to create public shares", func() {
        BeforeEeach(func() {
            // set up the mock
            // this is implicitly creating the expectation that it will be called Once()
            // this will throw an error if the method is not called
            gatewayClient.
                EXPECT().
                CheckPermission(
                    mock.Anything,
                    mock.Anything,
                ).
                Return(checkPermissionResponse, nil)
        })
        It("should return a public share", func() {
            // call the method
            req := &link.CreatePublicShareRequest{
                ResourceInfo: &providerpb.ResourceInfo{
                    Owner: &userpb.UserId{
                        OpaqueId: "alice",
                    },
                    Path: "./NewFolder/file.txt",
                },
                Grant: &link.Grant{
                    Permissions: &link.PublicSharePermissions{
                        Permissions: linkPermissions,
                    },
                    Password: "SecretPassw0rd!",
                },
                Description: "test",
            }
            res, err := provider.CreatePublicShare(ctx, req)
            Expect(err).ToNot(HaveOccurred())
            Expect(res.GetStatus().GetCode()).To(Equal(rpc.Code_CODE_OK))
            Expect(res.GetShare()).To(Equal(createdLink))
        })
    })
})
```

{{% hint type="tip" title="Mocking in oCIS" %}}
Use the constructor on new mocks to register the `AssertExpectations` method to be called at the end of the tests via the `t.Cleanup()` method.
{{% /hint %}}
