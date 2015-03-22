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

package control

import (
	"encoding/json"

	"github.com/flowcker/flowcker/atom"
	"github.com/flowcker/flowcker/atom/net"
	fc "github.com/flowcker/flowcker/common"
)

// LaunchControlAtom launches a control atom
func LaunchControlAtom(molecule *fc.Molecule, addr string, input chan fc.IPInbound) *atom.Atom {
	index := findControlAtom(molecule)
	if index == -1 {
		molecule.Atoms = append(molecule.Atoms, fc.Atom{ID: 0, Element: "control", Config: &json.RawMessage{'{', '}'}, Addr: ""})
		AddControlAtomLinks(molecule)
		index = findControlAtom(molecule)
	}

	// Add control atom address to molecule
	molecule.Atoms[index].Addr = addr

	// Create the configuration for the control atom, that is the whole molecule
	//molecule.Atoms[index].Config = &json.RawMessage{'{', '}'}
	b, _ := json.Marshal(&molecule)
	// TODO check error
	molecule.Atoms[index].Config.UnmarshalJSON(b)

	// Start Control atom, note that it does not use the standard or any middleware
	atom := atom.NewAtom(element, input)
	atom.Run()

	// Create Packet dispatcher for atom
	packetDispatcherConf := net.NewPacketDispatcherConf(0, molecule.Atoms, molecule.Links)
	net.NewPacketDispatcher(packetDispatcherConf, atom.GetOutput())

	// Pass the configuration to the control atom and start it
	atom.Setup(&molecule.Atoms[index])
	atom.StartElement()

	return atom
}

func findControlAtom(molecule *fc.Molecule) int {
	for index, atom := range molecule.Atoms {
		if atom.ID == 0 {
			return index
		}
	}

	return -1
}

// LaunchTCP launches control atom with TCP PortsServer
func LaunchTCP(molecule *fc.Molecule, addr string) (*atom.Atom, *net.PortsServer) {
	if addr == "" {
		addr = "0.0.0.0:0"
	}

	portsServer := net.NewPortsServer(addr)
	err := portsServer.Listen()
	if err != nil {
		panic(err)
	}

	atom := LaunchControlAtom(molecule, portsServer.GetAddr().String(), net.InboundUnirouter(portsServer.GetOutput()))
	return atom, portsServer
}
