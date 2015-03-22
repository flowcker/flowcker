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

package runtime

import (
	"os"

	fc "github.com/flowcker/flowcker/common"
	control "github.com/flowcker/flowcker/runtime/control"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("flowcker_runtime")

type AtomLauncher func(atom *fc.Atom) (string, error)

func Run(molecule *fc.Molecule, atomLauncher AtomLauncher) {
	for i := range molecule.Atoms {
		// Launch atom and find addr
		molecule.Atoms[i].Addr, _ = atomLauncher(&molecule.Atoms[i])
	}

	controlIP := os.Getenv("FLOWCKER_CONTROL_IP")
	if controlIP == "" {
		log.Warning("FLOWCKER_CONTROL_IP is not set, using default 127.0.0.1")
		controlIP = "127.0.0.1"
	}

	atom, _ := control.LaunchTCP(molecule, controlIP+":0")

	<-atom.CloseChannel()
}
