---
title: "Unit Testing"
date: 2024-04-25T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development/unit-testing
geekdocFilePath: _index.md
---

{{< toc >}}

Go is a statically typed language, which makes it easy to write unit tests. The Go standard library provides a `testing` package that allows you to write tests for your code. The testing package provides a framework for writing tests, and the `go test` command runs the tests. Other than that there are a lot of libraries and tools available to make testing easier.

- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks that plays nicely with the standard library.
- [Ginkgo](https://onsi.github.io/ginkgo/) - A BDD-style testing framework for Go.
- [Gomega](https://onsi.github.io/gomega/) - A matcher/assertion library for Ginkgo.
- [GoDog](https://github.com/cucumber/godog) - A Behavior-Driven Development framework for Go which uses Gherkin.

In oCIS we generally use [Ginkgo](https://onsi.github.io/ginkgo/) framework for testing. To keep things consistent, we would encourage you to use the same. In some cases, where you feel the need for a more verbose or more "code oriented" approach, you can also use the testing package from the standard library without ginkgo.

## 1 Ginkgo

Using a framework like [Ginkgo](https://onsi.github.io/ginkgo/) brings many advantages.

### Pros

- Provides a BDD-style syntax which makes it easier to write reusable and understandable tests
- Together with [Gomega](https://onsi.github.io/gomega/) it provides a powerful and expressive framework with assertions in a natural language
- Natural Language Format empowers testing in a way that resembles user interactions with the system
- In the context of microservices it is particularly well suited to test individual services and the interactions between them
- Offers support for asynchronous testing which makes it easier to test code that involves concurrency
- Nested and structured containers and setup capabilities make it easy to organize tests and adhere to the DRY principle
- Provides helpful error messages to identify and fix issues
- Very usable for Test Driven Development following the ["Red, Green, Cleanup, Repeat"](https://en.wikipedia.org/wiki/Test-driven_development) workflow.

### Cons

- Sometimes it can be difficult to get started with
- Asynchronous behaviour brings more complexity to tests.
- Not compatible with broadly known `testify` package

### Example

As you can see, **Ginkgo** and **Gomega** together provide the foundation to write understandable and maintainable tests which can mimic user interaction and the interactions between microservices.

```go
Describe("Public Share Provider", func() {
  Context("When the user has no share permission", func() {
    BeforeEach(func() {
        // downgrade user permissions to have no share permission
        resourcePermissions.AddGrant = false
    })
    It("should return grpc invalid argument", func() {
        req := &link.CreatePublicShareRequest{}

        res, err := provider.CreatePublicShare(ctx, req)
        Expect(err).ToNot(HaveOccurred())
        Expect(res.GetStatus().GetCode()).To(Equal(rpc.Code_CODE_INVALID_ARGUMENT))
        Expect(res.GetStatus().GetMessage()).To(Equal("no share permission"))
    })
})
```

### How to use it in oCIS

{{< button relref="testing-ginkgo" size="large" >}}{{< icon "gdoc_arrow_right_alt" >}} Read more{{< /button >}}

## 2 Testing Package

For smaller straight-forward tests of some packages it might feel more natural to use the testing package that comes with the go standard library.

### Pros

- Straightforward approach
- Naming conventions
- Built-in tooling

### Cons

- Difficult to reuse code in larger and more complex packages
- Difficult to create clean and isolated setups for the test steps
- No natural language resemblance


### How to use it in ocis

{{< button relref="testing-pkg" size="large" >}}{{< icon "gdoc_arrow_right_alt" >}} Read more{{< /button >}}
