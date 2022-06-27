package gncp

import (
	"errors"
	"net"
)

type CpConn struct {
	net.Conn
	pool *GncpPool
}

// Destroy will close connection and release connection from connection pool.
func (conn *CpConn) Destroy() error {
	if conn.pool == nil {
		return errors.New("Connection not belong any connection pool.")
	}
	err := conn.pool.Remove(conn.Conn)
	if err != nil {
		return err
	}
	conn.pool = nil
	return nil
}

// Close will push connection back to connection pool. It will not close the real connection.
func (conn *CpConn) Close() error {
	if conn.pool == nil {
		return errors.New("Connection not belong any connection pool.")
	}
	return conn.pool.Put(conn.Conn)
}
