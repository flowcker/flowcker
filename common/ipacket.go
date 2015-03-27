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

package common

import "fmt"

// IPacket represents an information packet
type IPacket struct {
	Type IPType
	Data []byte
}

// GetType returns IP type
func (p IPacket) GetType() IPType {
	return p.Type
}

// GetData returns IP data
func (p IPacket) GetData() []byte {
	return p.Data
}

// IPacketInbound contains information packets RECEIVED by atoms
type iPacketInbound struct {
	IP
	To Port
}

// GetTo returns the port the IP is coming into
func (p iPacketInbound) GetTo() Port {
	return p.To
}

// IPacketOutbound contains information packets SENT by atoms
type iPacketOutbound struct {
	IP
	From          Port
	IndexSelected bool
}

// GetFrom returns the port the IP is leaving from
func (p iPacketOutbound) GetFrom() Port {
	return p.From
}

// GetIndexSelected returns whether the outbound IP has a selected port index
func (p iPacketOutbound) GetIndexSelected() bool {
	return p.IndexSelected
}

func (p iPacketOutbound) GetAll() bool {
	return false
}

type iPacketOutboundAll struct {
	IP
	From Port
}

// GetFrom returns the port the IP is leaving from
func (p iPacketOutboundAll) GetFrom() Port {
	return p.From
}

// GetIndexSelected returns whether the outbound IP has a selected port index
func (p iPacketOutboundAll) GetIndexSelected() bool {
	return false
}

func (p iPacketOutboundAll) GetAll() bool {
	return true
}

// NewIP creates a DATA IP
func NewIP(data []byte) IP {
	return NewIPWithType(data, DATA)
}

// NewDisconnectIP creates a DISCONNECT IP
func NewDisconnectIP() IP {
	return NewIPWithType([]byte{}, DISCONNECT)
}

// NewIPWithType creates an IP with arbitrary type
func NewIPWithType(data []byte, ipType IPType) IP {
	return &IPacket{Type: ipType, Data: data}
}

// NewIPOut creates a IP to send outgoing
func NewIPOut(packet interface{}, port interface{}) IPOutbound {
	var port_ Port
	var indexSelected bool

	switch p := port.(type) {
	case string:
		port_.Name = p
		indexSelected = false
	case Port:
		port_ = p
		indexSelected = true
	default:
		panic(fmt.Sprintf("NewIPOut: unexpected type %T", p))
	}

	switch p := packet.(type) {
	case []byte:
		return &iPacketOutbound{IP: NewIP(p), From: port_, IndexSelected: indexSelected}
	case IP:
		return &iPacketOutbound{IP: p, From: port_, IndexSelected: indexSelected}
	default:
		panic(fmt.Sprintf("NewIPOut: unexpected type %T", p))
	}
}

// NewIPOutAll creates a IP to send outgoing to all index in the same port
func NewIPOutAll(packet interface{}, port string) IPOutbound {
	switch p := packet.(type) {
	case []byte:
		return &iPacketOutboundAll{IP: NewIP(p), From: Port{Name: port}}
	case IP:
		return &iPacketOutboundAll{IP: p, From: Port{Name: port}}
	default:
		panic(fmt.Sprintf("unexpected type %T", p))
	}
}

// NewIPIn creates a IP that is recieved
func NewIPIn(packet interface{}, port Port) IPInbound {
	switch p := packet.(type) {
	case []byte:
		return &iPacketInbound{NewIP(p), port}
	case string:
		return &iPacketInbound{NewIP([]byte(p)), port}
	case IP:
		return &iPacketInbound{p, port}
	default:
		panic(fmt.Sprintf("unexpected type %T", p))
	}
}

// Check interfaces
var _ IP = (*IPacket)(nil)
var _ IPInbound = (*iPacketInbound)(nil)
var _ IPOutbound = (*iPacketOutbound)(nil)
