# Rad [![GoDoc](https://godoc.org/github.com/purehyperbole/rad?status.svg)](https://godoc.org/github.com/purehyperbole/rad) [![Go Report Card](https://goreportcard.com/badge/github.com/purehyperbole/rad)](https://goreportcard.com/report/github.com/purehyperbole/rad) [![Build Status](https://travis-ci.org/purehyperbole/rad.svg?branch=master)](https://travis-ci.org/purehyperbole/rad)

A concurrent lock free radix tree implementation for go.

Rad is threadsafe and can handle concurrent reads/writes from any number of threads.

# Installation

To start using rad, you can run:

`$ go get github.com/purehyperbole/rad`

# Usage

To create a new radix tree

```go
package main

import (
    "github.com/purehyperbole/rad"
)

func main() {
    // create a new radix tree
    r := rad.New()
}
```

`Lookup` can be used to retrieve a stored value

```go
value := r.Lookup([]byte("myKey1234"))
```

`Insert` allows a value to be stored for a given key.

A successful insert will return true.

If the operation conflicts with an insert from another thread, it will return false.

Values must implement the `Comparable` interface to facilitate comparing one value to another.

There are some default implementations for `[]byte` and `string`. Please see the documentation for usage.

```go
if r.Insert([]byte("key"), &Thing{12345}) {
    fmt.Println("success!")
} else {
    fmt.Println("insert failed")
}
```

`MustInsert` can be used if you want to retry the insertion until it is successful.
```go
r.MustInsert([]byte("key"), &Thing{12345})
```

`Swap` can be used to atomically swap a value. It will return false if the swap was unsuccessful.
```go
r.Swap([]byte("my-key"), &Thing{"old-value"}, &Thing{"new-value"})
```

`Iterate` allows for iterating keys in the tree

```go
// iterate over all keys
r.Iterate(nil, func(key []byte, value interface{}) error {
    ...
    // if the returned error is not nil, the iterator will stop
    return nil
})

// iterate over all subkeys of "rad"
r.Iterate([]byte("rad"), func(key []byte, value interface{}) error {
    ...
    return nil
})
```

# Features/Wishlist

- [x] Lock free Insert using CAS (compare & swap)
- [x] Lookup
- [x] Basic key iterator
- [ ] Delete

## Why?

This project was created to learn about lock free data structures. As such, it probably should not be used for any real work. Use at your own risk!

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2018 purehyperbole.

Code released under
[the MIT License](LICENSE).
