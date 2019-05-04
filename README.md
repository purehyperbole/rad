# rad [![GoDoc](https://godoc.org/github.com/purehyperbole/rad?status.svg)](https://godoc.org/github.com/purehyperbole/rad) [![Go Report Card](https://goreportcard.com/badge/github.com/purehyperbole/rad)](https://goreportcard.com/report/github.com/purehyperbole/rad) [![Build Status](https://travis-ci.org/purehyperbole/rad.svg?branch=master)](https://travis-ci.org/purehyperbole/rad)

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

`Get` allows data to be retrieved.

```go
data, err := db.Get([]byte("myKey1234"))
```

`Set` allows data to be stored.

```go
err := db.Set([]byte("myKey1234"), []byte(`{"status": "ok"}`))
```

# Features/Wishlist

- [x] Persistence
- [x] Compressed tree nodes (radix)
- [x] Sync mmap data on resize
- [ ] Configurable sync on write options
- [ ] Transactions (MVCC)
- [ ] Data file compaction

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2018 purehyperbole.

Code released under
[the MIT License](LICENSE).
