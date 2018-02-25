# Useful links

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

# Running Man-In-The-Middle proxy

```
docker cp mysqlx_mysql_1:/var/lib/mysql/server-cert.pem internal/
docker cp mysqlx_mysql_1:/var/lib/mysql/server-key.pem internal/
```
