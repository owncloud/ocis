[![Travis](https://travis-ci.org/spacewander/go-suffix-tree.svg?branch=master)](https://travis-ci.org/spacewander/go-suffix-tree)
[![GoReportCard](http://goreportcard.com/badge/spacewander/go-suffix-tree)](http://goreportcard.com/report/spacewander/go-suffix-tree)
[![codecov.io](https://codecov.io/github/spacewander/go-suffix-tree/coverage.svg?branch=master)](https://codecov.io/github/spacewander/go-suffix-tree?branch=master)
[![license](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/spacewander/go-suffix-tree/blob/master/LICENSE)
[![godoc](https://img.shields.io/badge/godoc-reference-green.svg)](https://godoc.org/github.com/spacewander/go-suffix-tree)

# go-suffix-tree

This "suffix" package implements a [suffix tree](https://en.wikipedia.org/wiki/Suffix_tree).

As a suffix tree, it allows to lookup a key in O(k) operations.
In some cases(for example, some scenes in our production), this can be faster than a hash table because
   the hash function is an O(n) operation, with poor cache locality.

Plus suffix tree is more memory-effective than a hash table.

## Example

A simple use case:
```go
import (
    suffix "github.com/spacewander/go-suffix-tree"
)

var (
    TubeNameTree *suffix.Tree
    TubeNames = []string{
        // ...
    }
)

func init() {
    tree := suffix.NewTree()
    for _, s := range TubeNames {
        tree.Insert([]byte(s), &s)
    }
    TubeNameTree = tree
}

func getTubeName(name []byte) *string {
    res, found := TubeNameTree.Get(name)
    if found {
        return res.(*string)
    }
    return nil
}
```

For more usage, see the [godoc](https://godoc.org/github.com/spacewander/go-suffix-tree).
