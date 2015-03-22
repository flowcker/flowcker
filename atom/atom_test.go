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

import (
	"os"
	"sync"
	"testing"

	fc "github.com/flowcker/flowcker/common"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)

func setupTestLogging() {
	var logging_backend = logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(
			"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
		),
	)

	logging.SetBackend(logging_backend)
}

func TestIdentityAtom(t *testing.T) {
	setupTestLogging()
	input := make(chan fc.IPInbound)
	atom := NewAtom(identityElementForTests, input)
	output := atom.GetOutput()

	atom.Run()
	atom.Setup(nil)
	atom.StartElement()

	var barrier sync.WaitGroup

	barrier.Add(1)
	go func() {
		var outgoingP fc.IPOutbound

		outgoingP = <-output
		assert.Equal(t, outgoingP.GetData(), []byte{1, 2, 3})

		outgoingP = <-output
		assert.Equal(t, outgoingP.GetData(), []byte{4, 1, 2, 3})

		<-atom.CloseChannel()

		barrier.Done()
	}()

	assert.Equal(t, true, true)

	input <- fc.NewIPIn([]byte{1, 2, 3}, fc.Port{Name: "input", Index: 0})
	input <- fc.NewIPIn([]byte{4, 1, 2, 3}, fc.Port{Name: "input", Index: 0})
	close(input)

	barrier.Wait()
}

func TestIdentityAtomNoStart(t *testing.T) {
	setupTestLogging()
	input := make(chan fc.IPInbound)
	atom := NewAtom(identityElementForTests, input)
	output := atom.GetOutput()

	atom.Run()
	atom.Setup(nil)

	var barrier sync.WaitGroup

	barrier.Add(1)
	go func() {
		var outgoingP fc.IPOutbound

		outgoingP = <-output
		assert.Equal(t, outgoingP.GetData(), []byte{1, 2, 3})

		outgoingP = <-output
		assert.Equal(t, outgoingP.GetData(), []byte{4, 1, 2, 3})

		barrier.Done()
	}()

	assert.Equal(t, true, true)

	input <- fc.NewIPIn([]byte{1, 2, 3}, fc.Port{"input", 0})
	input <- fc.NewIPIn([]byte{4, 1, 2, 3}, fc.Port{"input", 0})

	barrier.Wait()
}

func identityElementForTests(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	out = make(chan fc.IPOutbound)

	go func() {
		defer log.Debug("Identity element: exiting")
		defer close(out)
		log.Debug("Starting Identity element")
		for incoming := range in {
			switch incoming.GetTo().Name {
			case "input":
				log.Debug("Identity element: received data")
				out <- fc.NewIPOut(incoming.GetData(), fc.Port{Name: "output"})
				log.Debug("Identity element: data sent")
			}
		}
	}()

	return out, nil
}
