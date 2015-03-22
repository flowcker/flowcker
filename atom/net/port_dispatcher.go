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
	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
)

// PortDispatcher dispaches packages from a specfic outport
type portDispatcher struct {
	dispatcher *atomDispatcher
	toAtom     uint32
	toPort     string
	toIndex    uint32
}

func newPortDispatcher(atomID uint32, portName string, portIndex uint32, atomDispatcher *atomDispatcher) (portD *portDispatcher) {
	portD = new(portDispatcher)

	portD.dispatcher = atomDispatcher
	portD.toAtom = atomID
	portD.toPort = portName
	portD.toIndex = portIndex

	return portD
}

func (portD *portDispatcher) dispatch(ip fc.IP) (err error) {
	log.Debug("portDispatcher: Dispatching IP with data \n%s", spew.Sdump(ip.GetData()))
	return portD.dispatcher.dispatch(newIPWireType(ip.GetType(), ip.GetData(), fc.Port{Name: portD.toPort, Index: portD.toIndex}, portD.toAtom))
}
