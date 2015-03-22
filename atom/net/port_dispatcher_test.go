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

	fc "github.com/flowcker/flowcker/common"
	"github.com/stretchr/testify/assert"
)

func TestPortDispatcherBasic(t *testing.T) {
	server := NewPortsServer("127.0.0.1:0")
	server.Listen()
	addr := server.GetAddr()
	time.Sleep(1000 * time.Millisecond)

	atomDispatcher := newAtomDispatcher(100, addr.String())
	dispatcher := newPortDispatcher(100, "testport", 0, atomDispatcher)
	dispatcher.dispatch(fc.NewIP([]byte{1, 2, 3}))

	p := <-server.GetOutput()
	assert.Equal(t, p.AtomID, uint32(100))
	assert.Equal(t, p.Port, fc.Port{Name: "testport", Index: 0})
	assert.Equal(t, p.Data, []byte{1, 2, 3})
}
