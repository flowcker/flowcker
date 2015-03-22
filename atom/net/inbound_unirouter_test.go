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

	fc "github.com/flowcker/flowcker/common"
	"github.com/stretchr/testify/assert"
)

func TestInboundUnirouterBasic(t *testing.T) {
	server := NewPortsServer("127.0.0.1:0")
	server.Listen()
	addr := server.GetAddr()

	out := InboundUnirouter(server.GetOutput())

	client := NewPortsClient(addr.String())
	client.Write(newIPWire([]byte{1, 2, 3}, fc.Port{"testport", 0}, 100))

	p := <-out
	assert.Equal(t, p.GetTo(), fc.Port{"testport", 0})
	assert.Equal(t, p.GetData(), []byte{1, 2, 3})
}
