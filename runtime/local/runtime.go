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

package local

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	fc "github.com/flowcker/flowcker/common"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("flowcker_runtime")

func LaunchAtom(atom *fc.Atom) (string, error) {
	args := strings.Split(atom.Element, " ")
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Env = []string{
		"HOST=127.0.0.1",
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(stdout)

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	addr, isPrefix, err := r.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	if isPrefix {
		log.Fatal("Address does not fit in buffer")
	}

	log.Debug("Atom %d launched", atom.ID)

	go func() {
		cmd.Wait()
	}()

	return string(addr), nil
}
