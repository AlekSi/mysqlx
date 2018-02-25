# mysqlx

MySQL driver for Go's (golang) `database/sql` package and MySQL X Protocol.

[![GoDoc](https://godoc.org/github.com/AlekSi/mysqlx?status.svg)](https://godoc.org/github.com/AlekSi/mysqlx)
[![Build Status](https://travis-ci.org/AlekSi/mysqlx.svg?branch=master)](https://travis-ci.org/AlekSi/mysqlx)
[![Codecov](https://codecov.io/gh/AlekSi/mysqlx/branch/master/graph/badge.svg)](https://codecov.io/gh/AlekSi/mysqlx)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlekSi/mysqlx)](https://goreportcard.com/report/github.com/AlekSi/mysqlx)

It requires Go 1.10+.

## Status

**Alpha quality. Do not use in production!**

You are, however, is encouraged to try it in development and report bugs.

## Data source format

```
mysqlx://username:password@host:port/database?_param=value&session_variable=value&…
```

All query parameters that are not starting with `_` are used as session variables
and are set whenever a connection is opened.
Parameters starting with `_` are listed below:

* `_auth-method`: `PLAIN` or `MYSQL41` (see [AuthMethod type](https://godoc.org/github.com/AlekSi/mysqlx#AuthMethod))

## TODO

* Real TLS support.
* Binary strings.
* Large uint64.
* More tests for correct connection closing.
* More concurrent tests.
* Benchmarks.
* Support for https://github.com/gogo/protobuf (?)
* Charsets.
* Time zones.
* Real prepared statements.
* Named values.
* Expose notices and warnings (?).
