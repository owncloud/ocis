// Copyright 2011 The Go Authors. All rights reserved.
// Copyright 2021 The LibreGraph Authors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldapserver

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"strings"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
	"github.com/libregraph/idm/pkg/ldapdn"
)

type Adder interface {
	Add(boundDN string, req *ldap.AddRequest, conn net.Conn) (LDAPResultCode, error)
}

type Binder interface {
	Bind(bindDN, bindSimplePw string, conn net.Conn) (LDAPResultCode, error)
}

type Deleter interface {
	Delete(boundDN string, req *ldap.DelRequest, conn net.Conn) (LDAPResultCode, error)
}

type Modifier interface {
	Modify(boundDN string, req *ldap.ModifyRequest, conn net.Conn) (LDAPResultCode, error)
}

type Searcher interface {
	Search(boundDN string, req *ldap.SearchRequest, conn net.Conn) (ServerSearchResult, error)
}

type Closer interface {
	Close(boundDN string, conn net.Conn) error
}

type Server struct {
	AddFns      map[string]Adder
	BindFns     map[string]Binder
	DeleteFns   map[string]Deleter
	ModifyFns   map[string]Modifier
	SearchFns   map[string]Searcher
	CloseFns    map[string]Closer
	Quit        chan bool
	EnforceLDAP bool
	Stats       *Stats
}

type ServerSearchResult struct {
	Entries    []*ldap.Entry
	Referrals  []string
	Controls   []ldap.Control
	ResultCode LDAPResultCode
}

func NewServer() *Server {
	s := new(Server)
	s.Quit = make(chan bool)

	d := defaultHandler{}
	s.AddFns = make(map[string]Adder)
	s.BindFns = make(map[string]Binder)
	s.DeleteFns = make(map[string]Deleter)
	s.ModifyFns = make(map[string]Modifier)
	s.SearchFns = make(map[string]Searcher)
	s.CloseFns = make(map[string]Closer)
	s.BindFunc("", d)
	s.SearchFunc("", d)
	s.CloseFunc("", d)
	s.Stats = nil
	return s
}

func (server *Server) AddFunc(baseDN string, f Adder) {
	server.AddFns[baseDN] = f
}

func (server *Server) BindFunc(baseDN string, f Binder) {
	server.BindFns[baseDN] = f
}

func (server *Server) DeleteFunc(baseDN string, f Deleter) {
	server.DeleteFns[baseDN] = f
}

func (server *Server) ModifyFunc(baseDN string, f Modifier) {
	server.ModifyFns[baseDN] = f
}

func (server *Server) SearchFunc(baseDN string, f Searcher) {
	server.SearchFns[baseDN] = f
}

func (server *Server) CloseFunc(baseDN string, f Closer) {
	server.CloseFns[baseDN] = f
}

func (server *Server) QuitChannel(quit chan bool) {
	server.Quit = quit
}

func (server *Server) ListenAndServeTLS(listenString string, certFile string, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}
	tlsConfig.ServerName = "localhost"
	ln, err := tls.Listen("tcp", listenString, &tlsConfig)
	if err != nil {
		return err
	}
	err = server.Serve(ln)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) SetStats(enable bool) {
	if enable {
		server.Stats = &Stats{}
	} else {
		server.Stats = nil
	}
}

func (server *Server) GetStats() Stats {
	return *server.Stats.Clone()
}

func (server *Server) ListenAndServe(listenString string) error {
	ln, err := net.Listen("tcp", listenString)
	if err != nil {
		return err
	}
	err = server.Serve(ln)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Serve(ln net.Listener) error {
	newConn := make(chan net.Conn)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if !strings.HasSuffix(err.Error(), "use of closed network connection") {
					log.Printf("Error accepting network connection: %s", err.Error())
				}
				break
			}
			log.Printf("New Connection on %s", ln.Addr())
			newConn <- conn
		}
	}()

listener:
	for {
		select {
		case c := <-newConn:
			server.Stats.countConns(1)
			go server.handleConnection(c)
		case <-server.Quit:
			ln.Close()
			break listener
		}
	}
	return nil
}

//
func (server *Server) handleConnection(conn net.Conn) {
	boundDN := "" // "" == anonymous

handler:
	for {
		// Read incoming LDAP packet.
		packet, err := ber.ReadPacket(conn)
		if err == io.EOF { // Client closed connection.
			break
		} else if err != nil {
			log.Printf("handleConnection ber.ReadPacket ERROR: %s", err.Error())
			break
		}

		// Sanity check this packet.
		if len(packet.Children) < 2 {
			log.Print("len(packet.Children) < 2")
			break
		}
		// Check the message ID and ClassType.
		messageID, ok := packet.Children[0].Value.(int64)
		if !ok {
			log.Print("malformed messageID")
			break
		}
		req := packet.Children[1]
		if req.ClassType != ber.ClassApplication {
			log.Print("req.ClassType != ber.ClassApplication")
			break
		}
		// Handle controls if present.
		controls := []ldap.Control{}
		if len(packet.Children) > 2 {
			for _, child := range packet.Children[2].Children {
				c, err := ldap.DecodeControl(child)
				if err != nil {
					log.Printf("handleConnection decode control ERROR: %s", err.Error())
					continue
				}
				controls = append(controls, c)
			}
		}

		// log.Printf("DEBUG: handling operation: %s [%d]", ldap.ApplicationMap[uint8(req.Tag)], req.Tag)
		// ber.PrintPacket(packet) // DEBUG

		// Dispatch the LDAP operation.
		switch req.Tag { // LDAP op code.
		default:
			responsePacket := encodeLDAPResponse(messageID, ldap.ApplicationAddResponse, ldap.LDAPResultOperationsError, "Unsupported operation: add")
			if err = sendPacket(conn, responsePacket); err != nil {
				log.Printf("sendPacket error %s", err.Error())
			}
			log.Printf("Unhandled operation: %s [%d]", ldap.ApplicationMap[uint8(req.Tag)], req.Tag)
			break handler

		case ldap.ApplicationAddRequest:
			server.Stats.countAdds(1)

			resultCode := uint16(ldap.LDAPResultSuccess)
			resultMsg := ""
			if err = HandleAddRequest(req, boundDN, server, conn); err != nil {
				var lErr *ldap.Error
				if errors.As(err, &lErr) {
					resultCode = lErr.ResultCode
					if lErr.Err != nil {
						resultMsg = lErr.Err.Error()
					}
				} else {
					resultCode = ldap.LDAPResultOperationsError
					resultMsg = err.Error()
				}
			}

			responsePacket := encodeLDAPResponse(messageID, ldap.ApplicationAddResponse, LDAPResultCode(resultCode), resultMsg)
			if err = sendPacket(conn, responsePacket); err != nil {
				log.Printf("sendPacket error %s", err.Error())
				break handler
			}

		case ldap.ApplicationBindRequest:
			server.Stats.countBinds(1)
			ldapResultCode := HandleBindRequest(req, server.BindFns, conn)
			if ldapResultCode == ldap.LDAPResultSuccess {
				boundDN, ok = req.Children[1].Value.(string)
				if !ok {
					log.Print("Malformed Bind DN")
					break handler
				}
				if boundDN, err = ldapdn.ParseNormalize(boundDN); err != nil {
					log.Printf("Error normalizing Bind DN: %s", err)
					break handler
				}

			}
			responsePacket := encodeBindResponse(messageID, ldapResultCode)
			if err = sendPacket(conn, responsePacket); err != nil {
				log.Printf("sendPacket error %s", err.Error())
				break handler
			}

		case ldap.ApplicationDelRequest:
			server.Stats.countDeletes(1)
			resultCode := uint16(ldap.LDAPResultSuccess)
			resultMsg := ""
			if err = HandleDeleteRequest(req, boundDN, server, conn); err != nil {
				var lErr *ldap.Error
				if errors.As(err, &lErr) {
					resultCode = lErr.ResultCode
					if lErr.Err != nil {
						resultMsg = lErr.Err.Error()
					}
				} else {
					resultCode = ldap.LDAPResultOperationsError
					resultMsg = err.Error()
				}
			}

			responsePacket := encodeLDAPResponse(messageID, ldap.ApplicationDelResponse, LDAPResultCode(resultCode), resultMsg)
			if err = sendPacket(conn, responsePacket); err != nil {
				log.Printf("sendPacket error %s", err.Error())
				break handler
			}

		case ldap.ApplicationModifyRequest:
			server.Stats.countModifies(1)
			resultCode := uint16(ldap.LDAPResultSuccess)
			resultMsg := ""
			if err = HandleModifyRequest(req, boundDN, server, conn); err != nil {
				var lErr *ldap.Error
				if errors.As(err, &lErr) {
					resultCode = lErr.ResultCode
					if lErr.Err != nil {
						resultMsg = lErr.Err.Error()
					}
				} else {
					resultCode = ldap.LDAPResultOperationsError
					resultMsg = err.Error()
				}
			}
			responsePacket := encodeLDAPResponse(messageID, ldap.ApplicationModifyResponse, LDAPResultCode(resultCode), resultMsg)
			if err = sendPacket(conn, responsePacket); err != nil {
				log.Printf("sendPacket error %s", err.Error())
				break handler
			}

		case ldap.ApplicationSearchRequest:
			server.Stats.countSearches(1)
			if doneControls, err := HandleSearchRequest(req, &controls, messageID, boundDN, server, conn); err != nil {
				log.Printf("handleSearchRequest error %s", err.Error()) // TODO: make this more testable/better err handling - stop using log, stop using breaks?
				e := err.(*ldap.Error)
				if err = sendPacket(conn, encodeSearchDone(messageID, LDAPResultCode(e.ResultCode), doneControls)); err != nil {
					log.Printf("sendPacket error %s", err.Error())
					break handler
				}
				break handler
			} else {
				if err = sendPacket(conn, encodeSearchDone(messageID, ldap.LDAPResultSuccess, doneControls)); err != nil {
					log.Printf("sendPacket error %s", err.Error())
					break handler
				}
			}
		case ldap.ApplicationUnbindRequest:
			server.Stats.countUnbinds(1)
			break handler // Simply disconnect.

		}
	}

	for _, c := range server.CloseFns {
		c.Close(boundDN, conn)
	}

	conn.Close()
	server.Stats.countConnsClose(1)
}

func sendPacket(conn net.Conn, packet *ber.Packet) error {
	_, err := conn.Write(packet.Bytes())
	if err != nil {
		log.Printf("Error Sending Message: %s", err.Error())
		return err
	}
	return nil
}

func routeFunc(dn string, funcNames []string) string {
	bestPick := ""
	for _, fn := range funcNames {
		if strings.HasSuffix(dn, fn) {
			l := len(strings.Split(bestPick, ","))
			if bestPick == "" {
				l = 0
			}
			if len(strings.Split(fn, ",")) > l {
				bestPick = fn
			}
		}
	}
	return bestPick
}

func encodeLDAPResponse(messageID int64, responseType uint8, ldapResultCode LDAPResultCode, message string) *ber.Packet {
	responsePacket := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
	responsePacket.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, messageID, "Message ID"))
	response := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ber.Tag(responseType), nil, ldap.ApplicationMap[responseType])
	response.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagEnumerated, uint64(ldapResultCode), "resultCode: "))
	response.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "matchedDN: "))
	response.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, message, "errorMessage: "))
	responsePacket.AppendChild(response)
	return responsePacket
}

type defaultHandler struct {
}

func (h defaultHandler) Bind(bindDN, bindSimplePw string, conn net.Conn) (LDAPResultCode, error) {
	return ldap.LDAPResultInvalidCredentials, nil
}

func (h defaultHandler) Search(boundDN string, req *ldap.SearchRequest, conn net.Conn) (ServerSearchResult, error) {
	return ServerSearchResult{make([]*ldap.Entry, 0), []string{}, []ldap.Control{}, ldap.LDAPResultSuccess}, nil
}

func (h defaultHandler) Close(boundDN string, conn net.Conn) error {
	conn.Close()
	return nil
}
