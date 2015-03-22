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

package atom

import "encoding/json"
import fc "github.com/flowcker/flowcker/common"

// ControlMessage is sent to/from ControlAtom
type ControlMessage struct {
	Type string
	Data []byte
}

// AtomSetup contains the information necessary to configure an Atom
//
// e.g.: Addresses of other atoms and links, so the PacketDispatcher can be configured
type AtomSetup struct {
	AtomID uint32
	fc.Molecule
}

// JSONDecode decodes a ControlMessage
func (msg *ControlMessage) JSONDecode(data []byte) {
	json.Unmarshal(data, msg)
}
