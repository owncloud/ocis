package ldap

// LDAP automatic reconnection mechanism, inspired by:
// https://gist.github.com/emsearcy/cba3295d1a06d4c432ab4f6173b65e4f#file-ldap_snippet-go

import (
	"errors"

	"github.com/go-ldap/ldap/v3"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

type ldapConnection struct {
	Conn  *ldap.Conn
	Error error
}

type ConnWithReconnect struct {
	conn    chan ldapConnection
	reset   chan *ldap.Conn
	retries int
	logger  *log.Logger
}

func NewLDAPWithReconnect(logger *log.Logger, ldapURI, bindDN, bindPassword string) ConnWithReconnect {
	conn := ConnWithReconnect{
		conn:    make(chan ldapConnection),
		reset:   make(chan *ldap.Conn),
		retries: 1,
		logger:  logger,
	}
	go conn.ldapAutoConnect(ldapURI, bindDN, bindPassword)
	return conn
}

func (c ConnWithReconnect) Search(sr *ldap.SearchRequest) (*ldap.SearchResult, error) {
	conn, err := c.GetConnection()

	if err != nil {
		return nil, err
	}

	var res *ldap.SearchResult
	for try := 0; try <= c.retries; try++ {
		res, err = conn.Search(sr)
		if !ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
			// non network error, return it to the client
			return res, err
		}

		c.logger.Debug().Msgf("Network Error. attempt %d", try)
		conn, err = c.reconnect(conn)
		if err != nil {
			return nil, err
		}
		c.logger.Debug().Msg("retrying LDAP Search")
	}
	// if we get here we reached the maximum retries. So return an error
	return nil, ldap.NewError(ldap.ErrorNetwork, errors.New("max retries"))
}

func (c ConnWithReconnect) GetConnection() (*ldap.Conn, error) {
	conn := <-c.conn
	if conn.Conn != nil && !ldap.IsErrorWithCode(conn.Error, ldap.ErrorNetwork) {
		c.logger.Debug().Msg("using existing Connection")
		return conn.Conn, conn.Error
	}
	return c.reconnect(conn.Conn)
}

func (c ConnWithReconnect) ldapAutoConnect(ldapURI, bindDN, bindPassword string) {
	l, err := c.ldapConnect(ldapURI, bindDN, bindPassword)
	if err != nil {
		c.logger.Error().Err(err).Msg("autoconnect could not get ldap Connection")
	}

	for {
		select {
		case resConn := <-c.reset:
			// Only close the connection and reconnect if the current
			// connection, matches the one we got via the reset channel.
			// If they differ we already reconnected
			if l != nil && l == resConn {
				c.logger.Debug().Msgf("closing connection %v", &l)
				l.Close()
			}
			if l == resConn || l == nil {
				c.logger.Debug().Msg("reconnecting to LDAP")
				l, err = c.ldapConnect(ldapURI, bindDN, bindPassword)
			} else {
				c.logger.Debug().Msg("already reconnected")
			}
		case c.conn <- ldapConnection{l, err}:
		}
	}
}

func (c ConnWithReconnect) ldapConnect(ldapURI, bindDN, bindPassword string) (*ldap.Conn, error) {
	c.logger.Debug().Msgf("Connecting to %s", ldapURI)
	l, err := ldap.DialURL(ldapURI)
	if err != nil {
		c.logger.Error().Err(err).Msg("could not get ldap Connection")
	} else {
		c.logger.Debug().Msg("LDAP Connected")
		if bindDN != "" {
			c.logger.Debug().Msgf("Binding as %s", bindDN)
			err = l.Bind(bindDN, bindPassword)
			if err != nil {
				c.logger.Error().Err(err).Msg("Bind failed")
				l.Close()
				return nil, err
			}

		}
	}

	return l, err
}

func (c ConnWithReconnect) reconnect(resetConn *ldap.Conn) (*ldap.Conn, error) {
	c.logger.Debug().Msg("LDAP connection reset")
	c.reset <- resetConn
	c.logger.Debug().Msg("Waiting for new connection")
	result := <-c.conn
	return result.Conn, result.Error
}
