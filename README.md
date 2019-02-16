# gjson - Fast JSON Parser written in pure Go

[![Build Status](https://travis-ci.org/yagi5/gjson.svg?branch=master)](https://travis-ci.org/yagi5/gjson)
[![Coverage Status](https://coveralls.io/repos/github/yagi5/gjson/badge.svg?branch=master)](https://coveralls.io/github/yagi5/gjson?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/yagi5/gjson)](https://goreportcard.com/report/github.com/yagi5/gjson)
[![GoDoc](https://godoc.org/github.com/yagi5/gjson?status.svg)](https://godoc.org/github.com/yagi5/gjson)

### Installation

```shell
$ go get -u github.com/yagi5/gjson
```

### Usage

```go
package main

import (
	"fmt"

	"github.com/yagi5/gjson"
)

func main() {
	js := `{
			"key": "a", 
			"key2": [1, "a", true, null]
		}`
	jsMap, err := gjson.Decode([]byte(js))
	if err != nil {
		panic(err)
	}
	for key, val := range jsMap {
		fmt.Println(key)
		fmt.Println(val)
	}
}

// Output
// key
// a
// key2
// [1 a true nil]
```

### Lisence

MIT
