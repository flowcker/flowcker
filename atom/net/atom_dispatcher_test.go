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

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
	"github.com/stretchr/testify/assert"
)

func TestAtomDispatcherBasic(t *testing.T) {
	server := NewPortsServer("127.0.0.1:0")
	server.Listen()
	addr := server.GetAddr()
	time.Sleep(1000 * time.Millisecond)

	dispatcher := newAtomDispatcher(100, addr.String())
	dispatcher.dispatch(newIPWire([]byte{1, 2, 3}, fc.Port{Name: "testport", Index: 0}, 100))

	p := <-server.GetOutput()
	spew.Dump(p)
	spew.Dump(p.AtomID)
	assert.Equal(t, uint32(100), p.AtomID, "Equal atom id")
	assert.Equal(t, p.Port, fc.Port{"testport", 0})
	assert.Equal(t, p.Data, []byte{1, 2, 3})
}
