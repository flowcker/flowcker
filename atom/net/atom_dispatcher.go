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

// AtomDispatcher dispatches packages to a specific Atom
type atomDispatcher struct {
	toAtom uint32
	client *PortsClient
}

// NewAtomDispatcher creates a new AtomDispatcher
func newAtomDispatcher(atomID uint32, atomAddr string) (ad *atomDispatcher) {
	log.Debug("Creating atomDispatcher for atom %d", atomID)

	ad = new(atomDispatcher)
	ad.toAtom = atomID
	ad.client = NewPortsClient(atomAddr)

	return ad
}

func (ad *atomDispatcher) dispatch(p *iPacketWire) (err error) {
	if p.AtomID != ad.toAtom {
		log.Warning("AtomDispatcher sending ipacket with different AtomID")
	}

	log.Debug("atomDispatcher sending package to atom %d, port %s[%d] ", p.AtomID, p.Port.Name, p.Port.Index)

	return ad.client.Write(p)
}
