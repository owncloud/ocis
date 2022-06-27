# gncp [![Doc](http://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/eternnoir/gncp) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/eternnoir/gncp/blob/master/LICENSE) [![Build](https://travis-ci.org/eternnoir/gncp.svg?branch=develop)](https://travis-ci.org/eternnoir/gncp)
A thread safe connection pool for net.conn interface. Easy to manage, reuse and limit connections in golang.

## Install

Use `go get` to install package:

```
go get github.com/eternnoir/gncp
```

In source code:

```go
import "github.com/eternnoir/gncp"
```

## Usage

Full document: https://godoc.org/github.con/eternnoir/gncp

```go

// connCreator let connection know how to create new connection.
func connCreator() (net.Conn, error) {
	return net.Dial("tcp", "127.0.0.1:5566")
}

// Create new connection pool. It will initialize 3 connection in pool when pool created.
// If connection not enough in pool, pool will call creator to create new connection.
// But when total connection number pool created reach 10 connection, pool will not creat
// any new connection until someone call Remove().
pool, err := gncp.NewPool(3, 10, connCreator)

// Get connection from pool. If pool has no connection and total connection reach max number
// of connections, this method will block until someone put back connection to pool.
conn, err := pool.Get()

// Get connection from pool with timeout. It will wait one second, if still cannot get connection
// it will return timeout error.
conn, err := pool.GetWithTimeout(time.Duration(1) * time.Second)

// After you are finished using the connection call Close() method to put connection back to pool.
// It will not close real connection.
err := conn.Close()

// Remove connection from connection pool. The connection will not belong pool anymore.
// And this method will close connection.
err := pool.Remove(conn)

// Close connection pool. All connections in pool will be closed.
err := pool.Close()
```

## License

The MIT License (MIT)
