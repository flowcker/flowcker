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
	"strconv"

	fc "github.com/flowcker/flowcker/common"
)

// AddControlAtomLinks adds links between the control atom and all other atoms
func AddControlAtomLinks(molecule *fc.Molecule) {
	for _, atom := range molecule.Atoms {
		if atom.ID == 0 {
			continue
		}

		portName := "atom_" + strconv.FormatUint(uint64(atom.ID), 10)

		molecule.Links = append(molecule.Links, fc.Link{
			ID: 200, // TODO, do we need this ID ???
			From: fc.LinkSide{
				AtomID: 0,
				Port:   fc.Port{Name: portName, Index: 0},
			},
			To: fc.LinkSide{
				AtomID: atom.ID,
				Port:   fc.Port{Name: "control", Index: 0},
			},
		})

		molecule.Links = append(molecule.Links, fc.Link{
			ID: 201, // TODO, do we need this ID ???
			From: fc.LinkSide{
				AtomID: atom.ID,
				Port:   fc.Port{Name: "control", Index: 0},
			},
			To: fc.LinkSide{
				AtomID: 0,
				Port:   fc.Port{Name: portName, Index: 0},
			},
		})
	}
}
