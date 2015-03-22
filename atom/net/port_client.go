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
	"net"
)

// PortsClient connects to the PortServer on another atom
//
// PortsClient is used by PacketDispatcher to send information packets
type PortsClient struct {
	conn *net.TCPConn
	w    *bufio.Writer
}

// NewPortsClient creates a *PortClient
//
// The only parameter serverAddr is a string with the IP address and
// port of the atom PortServer. For example "192.168.0.16:5000".
func NewPortsClient(serverAddr string) (pc *PortsClient) {
	log.Debug("Creating PortsClient to %s", serverAddr)

	pc = new(PortsClient)

	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	pc.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	pc.w = bufio.NewWriter(pc.conn)

	return pc
}

// Write sends an information packet to the other side PortServer
func (pc *PortsClient) Write(packet *iPacketWire) (err error) {
	log.Debug("PortsClient sending packet")
	return packet.WriteTo(pc.w)
}
