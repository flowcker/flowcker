/*  Copyright (C) 2015 Leopoldo Lara Vazquez.

    This program is free software: you can redistribute it and/or  modify
    it under the terms of the GNU Affero General Public License, version 3,
	  as published by the Free Software Foundation.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package net

import (
	"bufio"
	"io"
	"net"
)

// PortsServer listen and receives information packets for an Atom
type PortsServer struct {
	laddr    *net.TCPAddr
	out      chan *iPacketWire
	listener *net.TCPListener
}

// NewPortsServer creates a PortsServer
func NewPortsServer(serverAddr string) (s *PortsServer) {
	log.Debug("Creating PortsServer binded to %s", serverAddr)

	s = new(PortsServer)
	s.laddr, _ = net.ResolveTCPAddr("tcp", serverAddr)
	s.out = make(chan *iPacketWire)

	return s
}

// Listen makes the PortsServer listen to new connections
func (s *PortsServer) Listen() (err error) {
	s.listener, err = net.ListenTCP("tcp", s.laddr)
	if err != nil {
		log.Debug("Error listening in PortsServer")
		return
	}
	log.Debug("Creating PortsServer listening to %s", s.GetAddr().String())

	go func() {
		for {
			err := s.accept()
			if err != nil {
				panic(err)
			}
		}
	}()

	return
}

// GetOutput returns a channel where new IPs arriving to this
// PortsServer will be send to
func (s *PortsServer) GetOutput() (out chan *iPacketWire) {
	return s.out
}

// GetAddr returns IP address the server is listening to
func (s *PortsServer) GetAddr() net.Addr {
	if s.listener == nil {
		log.Fatal("Accesing address when PortsServer is no listening")
	}
	return s.listener.Addr()
}

func (s *PortsServer) accept() (err error) {
	conn, err := s.listener.AcceptTCP()
	if err != nil {
		return err
	}

	go s.connectionHandler(conn)

	return
}

func (s *PortsServer) connectionHandler(conn *net.TCPConn) {
	log.Debug("PortsServer on %s, new client from %s", s.GetAddr().String(), conn.RemoteAddr().String())

	defer conn.Close()

	r := bufio.NewReader(conn)

	for {
		p := new(iPacketWire)
		err := p.ReadFrom(r)
		if err == io.EOF {
			log.Debug("PortsServer: Closed connection on %s, client from %s", s.GetAddr().String(), conn.RemoteAddr().String())
			break
		}
		if err != nil {
			panic(err)
		}

		log.Debug("PortsServer: Received packet on %s, from %s", s.GetAddr().String(), conn.RemoteAddr().String())
		s.out <- p
	}
}
