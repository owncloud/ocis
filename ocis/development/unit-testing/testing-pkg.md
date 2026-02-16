---
title: "Standard Library Testing"
date: 2024-04-25T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development/unit-testing
geekdocFilePath: testing-pkg.md

---

## Using the standard library

To write a unit test for your package, create a file with the `_test.go` suffix. For example, if you have a package `foo` with a file `foo.go`, you can create a file `foo_test.go` in the same directory. The test file should have the same package name as the package being tested. By doing this, you can access all exported and unexported identifiers of the package. It is a good practice to keep the test file in the same package as the code being tested.

### Simple Example

We are using an oversimplified example from [FooBarQuix](https://codingdojo.org/kata/FooBarQix/) to demonstrate how to use the `testing` package.

```go
package divide

import "strconv"

// If the number is divisible by 3, write "Yes" otherwise, the number
func IsDivisible(input int) string {
    if  (input % 3) == 0 {
        return "Yes"
    }
    return strconv.Itoa(input)
}
```

To test the `IsDivisible` function, create a file `divide_test.go` in the same directory as `divide.go`. The test file should have the same package name as the package being tested.

A test function in Go starts with `Test` and takes `*testing.T` as the only parameter. In most cases, you will name the unit test `Test[NameOfFunction]`. The testing package provides tools to interact with the test workflow, such as `t.Errorf`, which indicates that the test failed by displaying an error message on the console.

The test function for the `IsDivisible` function could look like this

```go
package divide

import "testing"

func TestDivide3(t *testing.T) {
    result := IsDivisible(3)
    if result != "Yes" {
        t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Yes")
    }
}
```

To run the test, use the `go test` command in the directory where the test file is located.

### Use a helper package for assertions

You could make the test more readable by using testify. The `assert` package provides a lot of helper functions to make the test more readable.

```go
package divide

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDivide3(t *testing.T) {
    result := IsDivisible(3)
    assert.Equal(t, "Yes", result)
}
```

### Table Driven Example

Write Table Driven Tests to test multiple inputs.

```go
package divide

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestIsDivisibleTableDriven(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name string
		input int
		want  string
	}{
		// the table itself
		{"9 should be Yes", 9, "Yes"},
		{"3 should be Yes", 3, "Yes"},
		{"1 is not Yes", 1, "1"},
		{"0 should be Yes", 0, "Yes"},
	}

	// The execution loop
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            answer := IsDivisible(tt.input)
            assert.Equal(t, tt.want, answer)
        })
    }
}
```

A table-driven test starts by defining the input structure. This can be seen like defining the columns of the table. Each row of the table lists a test case to execute. Once the table is defined, the execution loop can be created.

The execution loop calls `t.Run()`, which defines a subtest. In our example each row of the table defines a subtest named `[NameOfTheFuction]/[NameOfTheSubTest]`.

This way of writing tests is very popular, and considered the canonical way to write unit tests in Go.
