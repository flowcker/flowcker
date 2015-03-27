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

package control_test

import (
	"encoding/json"
	"sync"
	"testing"

	stdlib "github.com/flowcker/flowcker-stdlib/lib"
	"github.com/flowcker/flowcker/atom"
	fc "github.com/flowcker/flowcker/common"
	"github.com/flowcker/flowcker/runtime/control"
	"github.com/stretchr/testify/assert"
)

// TestBasicControl is basically an integration test
func TestBasicControl(t *testing.T) {
	setupTestLogging()

	var molecule fc.Molecule

	molecule.Atoms = append(molecule.Atoms, fc.Atom{ID: 1, Element: "identity", Config: &json.RawMessage{'{', '}'}, Addr: ""})
	molecule.Atoms = append(molecule.Atoms, fc.Atom{ID: 2, Element: "tester", Config: &json.RawMessage{'{', '}'}, Addr: ""})

	_, atomPortsServer := atom.LaunchTCP(stdlib.Identity, "127.0.0.1:0")
	molecule.Atoms[0].Addr = atomPortsServer.GetAddr().String()

	var barrier sync.WaitGroup
	barrier.Add(1)

	_, atomPortsServer = atom.LaunchTCP(func(config *fc.Atom, in chan fc.IPInbound) (chan fc.IPOutbound, error) {
		out := make(chan fc.IPOutbound)

		go func() {
			log.Debug("Tester launched")

			out <- fc.NewIPOut([]byte{1, 2, 3, 4}, fc.Port{Name: "output", Index: 0})

			incoming := <-in
			assert.Equal(t, []byte{1, 2, 3, 4}, incoming.GetData())

			barrier.Done()
		}()

		return out, nil
	}, "127.0.0.1:0")
	molecule.Atoms[1].Addr = atomPortsServer.GetAddr().String()

	molecule.Links = append(molecule.Links, fc.Link{
		ID: 400,
		From: fc.LinkSide{
			Port:   fc.Port{Name: "output", Index: 0},
			AtomID: 2,
		},
		To: fc.LinkSide{
			Port:   fc.Port{Name: "input", Index: 0},
			AtomID: 1,
		},
	})

	molecule.Links = append(molecule.Links, fc.Link{
		ID: 401,
		From: fc.LinkSide{
			Port:   fc.Port{Name: "output", Index: 0},
			AtomID: 1,
		},
		To: fc.LinkSide{
			Port:   fc.Port{Name: "input", Index: 0},
			AtomID: 2,
		},
	})

	_, _ = control.LaunchTCP(&molecule, "127.0.0.1:0")

	barrier.Wait()
}
