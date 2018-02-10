// mysqlx - MySQL driver for Go's database/​sql package and MySQL X Protocol.
// Copyright (c) 2017 Alexey Palazhchenko
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build ignore

package main

import (
	"crypto/tls"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"

	driver "github.com/AlekSi/mysqlx"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_connection"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_session"
	"github.com/AlekSi/mysqlx/internal/proto/mysqlx_sql"
)

func readClientMessage(r io.Reader) (proto.Message, []byte, error) {
	var head [5]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return nil, nil, err
	}

	buf := make([]byte, binary.LittleEndian.Uint32(head[:])+4)
	copy(buf, head[:])
	if _, err := io.ReadFull(r, buf[5:]); err != nil {
		return nil, nil, err
	}

	t := mysqlx.ClientMessages_Type(buf[4])
	var m proto.Message
	switch t {
	case mysqlx.ClientMessages_CON_CAPABILITIES_GET:
		m = new(mysqlx_connection.CapabilitiesGet)
	case mysqlx.ClientMessages_CON_CAPABILITIES_SET:
		m = new(mysqlx_connection.CapabilitiesSet)
	case mysqlx.ClientMessages_SESS_AUTHENTICATE_START:
		m = new(mysqlx_session.AuthenticateStart)
	case mysqlx.ClientMessages_SESS_AUTHENTICATE_CONTINUE:
		m = new(mysqlx_session.AuthenticateContinue)
	case mysqlx.ClientMessages_SQL_STMT_EXECUTE:
		m = new(mysqlx_sql.StmtExecute)
	default:
		return nil, nil, fmt.Errorf("readClientMessage: unhandled type of client message: %s (%d)", t, t)
	}

	if err := proto.Unmarshal(buf[5:], m); err != nil {
		return nil, buf, fmt.Errorf("readClientMessage: %s", err)
	}

	return m, buf, nil
}

func main() {
	log.SetFlags(0)

	listenF := flag.String("listen", "127.0.0.1:33061", "listen on that address")
	connectF := flag.String("connect", "127.0.0.1:33060", "connect to that address")
	serverCertF := flag.String("server-cert", "server-cert.pem", "server certificate")
	serverKeyF := flag.String("server-key", "server-key.pem", "server private key")
	flag.Parse()

	var tlsClientConfig *tls.Config
	if *serverCertF != "" || *serverKeyF != "" {
		cert, err := tls.LoadX509KeyPair(*serverCertF, *serverKeyF)
		if err != nil {
			log.Fatal(err)
		}
		tlsClientConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	log.Printf("Listening on %s...", *listenF)
	l, err := net.Listen("tcp", *listenF)
	if err != nil {
		log.Fatal(err)
	}

	client, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Accepted connection from %s to %s.", client.RemoteAddr(), client.LocalAddr())

	log.Printf("Connecting to %s...", *connectF)
	server, err := net.Dial("tcp", *connectF)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %s from %s.", server.RemoteAddr(), server.LocalAddr())

	startTLSClient := make(chan struct{})
	startTLSServer := make(chan struct{})

	// read from client, write to server
	go func() {
		prefix := fmt.Sprintf(
			"\n%s -> %s -> %s -> %s:\n",
			client.RemoteAddr(), client.LocalAddr(), server.LocalAddr(), server.RemoteAddr(),
		)
		logger := log.New(os.Stderr, prefix, log.Flags())
		for {
			m, b, err := readClientMessage(client)
			if err != nil {
				logger.Fatal(err)
			}

			msg := fmt.Sprintf("%T\n%s", m, m)
			switch auth := m.(type) {
			case *mysqlx_session.AuthenticateStart:
				msg += fmt.Sprintf("\nAuthData = %x", auth.AuthData)
			case *mysqlx_session.AuthenticateContinue:
				msg += fmt.Sprintf("\nAuthData = %x", auth.AuthData)
			case *mysqlx_session.AuthenticateOk:
				msg += fmt.Sprintf("\nAuthData = %x", auth.AuthData)
			}
			logger.Print(msg)

			var tlsRequested bool
			if set, ok := m.(*mysqlx_connection.CapabilitiesSet); ok {
				for _, cap := range set.GetCapabilities().GetCapabilities() {
					if cap.GetName() == "tls" && cap.GetValue().GetScalar().GetVBool() {
						tlsRequested = true
						break
					}
				}
			}

			if tlsRequested {
				close(startTLSClient)
			}

			_, err = server.Write(b)
			if err != nil {
				logger.Fatal(err)
			}

			if tlsRequested {
				<-startTLSServer
				startTLSServer = nil

				logger.Printf("Establishing TLS connection...")
				tlsClient := tls.Server(client, tlsClientConfig)
				if err = tlsClient.Handshake(); err != nil {
					logger.Panic(err)
				}
				client = tlsClient
				logger.Printf("TLS connection established.")
			}
		}
	}()

	// read from server, write to client
	go func() {
		prefix := fmt.Sprintf(
			"\n%s <- %s <- %s <- %s:\n",
			client.RemoteAddr(), client.LocalAddr(), server.LocalAddr(), server.RemoteAddr(),
		)
		logger := log.New(os.Stderr, prefix, log.Flags())
		for {
			m, b, err := driver.ReadMessage(server)
			if err != nil {
				logger.Fatal(err)
			}

			msg := fmt.Sprintf("%T\n%s", m, m)
			switch auth := m.(type) {
			case *mysqlx_session.AuthenticateStart:
				msg += fmt.Sprintf("\nAuthData = %x", auth.AuthData)
			case *mysqlx_session.AuthenticateContinue:
				msg += fmt.Sprintf("\nAuthData = %x", auth.AuthData)
			case *mysqlx_session.AuthenticateOk:
				msg += fmt.Sprintf("\nAuthData = %x", auth.AuthData)
			}
			logger.Print(msg)

			_, err = client.Write(b)
			if err != nil {
				logger.Fatal(err)
			}

			_, ok := m.(*mysqlx.Ok)
			if ok {
				select {
				case <-startTLSClient:
					close(startTLSServer)
					startTLSClient = nil

					logger.Printf("Establishing TLS connection...")
					tlsServer := tls.Client(server, &tls.Config{
						InsecureSkipVerify: true,
					})
					if err = tlsServer.Handshake(); err != nil {
						logger.Panic(err)
					}
					server = tlsServer
					logger.Printf("TLS connection established.")
				default:
				}
			}
		}
	}()

	select {}
}
