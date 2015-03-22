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

import fc "github.com/flowcker/flowcker/common"

// InboundUnirouter converts incoming IPacketRouted into IPacketInbound
func InboundUnirouter(in chan *iPacketWire) chan fc.IPInbound {
	out := make(chan fc.IPInbound)
	go func() {
		for incoming := range in {
			out <- fc.NewIPIn(incoming.Data, incoming.Port)
		}
	}()
	return out
}
