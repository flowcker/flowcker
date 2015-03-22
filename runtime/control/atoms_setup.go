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
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/flowcker/flowcker/atom"
	fc "github.com/flowcker/flowcker/common"
)

func getAtomSetupMessage(atomID uint32, molecule *fc.Molecule) (msg *atom.ControlMessage) {
	msg = new(atom.ControlMessage)
	msg.Type = "AtomSetup"

	var setup atom.AtomSetup

	setup.AtomID = atomID
	setup.Molecule = *molecule

	var err error
	msg.Data, err = json.Marshal(setup)
	if err != nil {
		spew.Dump(setup)
		panic("Error marshalling AtomSetup")
	}

	return msg
}

func sendAtomSetupMessage(atomID uint32, molecule *fc.Molecule, out chan fc.IPOutbound) {
	log.Info("Control Atom Element: sending setup to atom ID(%d)", atomID)

	data, err := json.Marshal(getAtomSetupMessage(atomID, molecule))
	if err != nil {
		spew.Dump(getAtomSetupMessage(atomID, molecule))
		panic(err)
	}
	out <- fc.NewIPOut(data, "atom_"+strconv.FormatUint(uint64(atomID), 10))

	log.Info("Control Atom Element: setup sent to atom ID(%d)", atomID)
}

func sendAtomsSetupMessage(molecule *fc.Molecule, out chan fc.IPOutbound) {
	log.Notice("Control Atom Element: sending setup to atoms")
	for _, atom := range molecule.Atoms {
		if atom.ID == 0 { // Skip control atom
			continue
		}
		sendAtomSetupMessage(atom.ID, molecule, out)
	}

	log.Notice("Control Atom Element: setup sent to all atoms")
}

func sendAtomsStartMessage(molecule *fc.Molecule, out chan fc.IPOutbound) {
	log.Notice("Control Atom Element: sending start packet to atoms")
	msg := new(atom.ControlMessage)
	msg.Type = "AtomStart"

	for _, atom := range molecule.Atoms {
		log.Info("Control Atom Element: sending start packet to atom ID(%d)", atom.ID)
		if atom.ID == 0 { // Skip control atom
			continue
		}

		data, _ := json.Marshal(msg)
		out <- fc.NewIPOut(data, "atom_"+strconv.FormatUint(uint64(atom.ID), 10))
	}
	log.Notice("Control Atom Element: start packet sent to atoms")
}
