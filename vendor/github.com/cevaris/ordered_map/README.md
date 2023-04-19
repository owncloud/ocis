# Ordered Map for golang

[![Build Status](https://travis-ci.org/cevaris/ordered_map.svg?branch=master)](https://travis-ci.org/cevaris/ordered_map)

**OrderedMap** is a Python port of OrderedDict implemented in golang. Golang's builtin `map` purposefully randomizes the iteration of stored key/values. **OrderedMap** struct preserves inserted key/value pairs; such that on iteration, key/value pairs are received in inserted (first in, first out) order.


## Features
- Full support Key/Value for all data types
- Exposes an Iterator that iterates in order of insertion
- Full Get/Set/Delete map interface
- Supports Golang v1.3 through v1.12

## Download and Install 
  
`go get https://github.com/cevaris/ordered_map.git`


## Examples

### Create, Get, Set, Delete

```go
package main

import (
    "fmt"
    "github.com/cevaris/ordered_map"
)

func main() {

    // Init new OrderedMap
    om := ordered_map.NewOrderedMap()

    // Set key
    om.Set("a", 1)
    om.Set("b", 2)
    om.Set("c", 3)
    om.Set("d", 4)

    // Same interface as builtin map
    if val, ok := om.Get("b"); ok == true {
        // Found key "b"
        fmt.Println(val)
    }

    // Delete a key
    om.Delete("c")

    // Failed Get lookup becase we deleted "c"
    if _, ok := om.Get("c"); ok == false {
        // Did not find key "c"
        fmt.Println("c not found")
    }
    
    fmt.Println(om)
}
```


### Iterator

```go
n := 100
om := ordered_map.NewOrderedMap()

for i := 0; i < n; i++ {
    // Insert data into OrderedMap
    om.Set(i, fmt.Sprintf("%d", i * i))
}

// Iterate though values
// - Values iteration are in insert order
// - Returned in a key/value pair struct
iter := om.IterFunc()
for kv, ok := iter(); ok; kv, ok = iter() {
    fmt.Println(kv, kv.Key, kv.Value)
}
```

### Custom Structs

```go
om := ordered_map.NewOrderedMap()
om.Set("one", &MyStruct{1, 1.1})
om.Set("two", &MyStruct{2, 2.2})
om.Set("three", &MyStruct{3, 3.3})

fmt.Println(om)
// Ouput: OrderedMap[one:&{1 1.1},  two:&{2 2.2},  three:&{3 3.3}, ]
```
  
## For Development

Git clone project 

`git clone https://github.com/cevaris/ordered_map.git`  
  
Build and install project

`make`

Run tests 

`make test`







