# mysqlx

MySQL driver for Go's (golang) `database/sql` package and MySQL X Protocol.

[![GoDoc](https://godoc.org/github.com/AlekSi/mysqlx?status.svg)](https://godoc.org/github.com/AlekSi/mysqlx)
[![Build Status](https://travis-ci.org/AlekSi/mysqlx.svg?branch=master)](https://travis-ci.org/AlekSi/mysqlx)
[![Codecov](https://codecov.io/gh/AlekSi/mysqlx/branch/master/graph/badge.svg)](https://codecov.io/gh/AlekSi/mysqlx)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlekSi/mysqlx)](https://goreportcard.com/report/github.com/AlekSi/mysqlx)

It requires Go 1.10+.

**Experimental work in progress, do not use in production.**

* https://dev.mysql.com/worklog/task/?id=8338
* https://dev.mysql.com/worklog/task/?id=8639
* https://dev.mysql.com/worklog/task/?id=9271
* https://dev.mysql.com/worklog/task/?id=10237
* https://dev.mysql.com/worklog/task/?id=10992
* https://dev.mysql.com/doc/internals/en/x-protocol.html
* https://dev.mysql.com/doc/dev/mysql-server/latest/PAGE_PROTOCOL.html
* https://dev.mysql.com/doc/refman/5.5/en/x-plugin.html
* https://dev.mysql.com/doc/refman/5.7/en/x-plugin.html
* https://dev.mysql.com/doc/refman/8.0/en/x-plugin.html

## Datasource parameters

* `_auth-method`

## TODO
* https://dev.mysql.com/doc/mysql-shell-excerpt/5.7/en/mysql-shell-connection-using-uri.html
* Binary strings.
* Large uint64.
* Correct connection closing.
* Concurrent tests.
* Benchmarks.
* Connection string format.
* TLS.
* Add support for https://github.com/gogo/protobuf
* Charsets.
* Time zones.
* Real prepared statements.
* Named values.
* Expose notices and warnings.
* MUCH MORE.

```
docker cp mysqlx_mysql_1:/var/lib/mysql/server-cert.pem internal/
docker cp mysqlx_mysql_1:/var/lib/mysql/server-key.pem internal/
```
