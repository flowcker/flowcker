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

// Command line interface for flowckerr-local
// Runs a molecule as local processes, only remember to put your libraries
// in the path
package main

import (
	"fmt"
	"os"

	fc "github.com/flowcker/flowcker/common"
	"github.com/flowcker/flowcker/runtime"
	"github.com/flowcker/flowcker/runtime/local"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Must pass molecule")
		return
	}
	runtime.Run(fc.NewMoleculeFromFile(os.Args[1]), local.LaunchAtom)
}
