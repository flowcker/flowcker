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
	"time"

	fc "github.com/flowcker/flowcker/common"
)

// Element is the core of controlling the other atoms
func element(config *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	out = make(chan fc.IPOutbound)

	var molecule fc.Molecule

	json.Unmarshal(*config.Config, &molecule)

	go func() {
		defer close(out)
		log.Debug("Starting Control Atom Element")

		// Setup all the atoms syncrhonously
		log.Debug("Control Atom Element: sending setup to atoms")
		sendAtomsSetupMessage(&molecule, out)
		// TODO make sure everyone recevie it? Wait time?
		time.Sleep(time.Duration(len(molecule.Atoms)) * 2 * time.Second)
		// Send start message to all atoms
		sendAtomsStartMessage(&molecule, out)
		log.Debug("Control Atom Element: setup and starting ended")
	}()

	return out, nil
}
