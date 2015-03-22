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

import "testing"
import "bytes"
import "bufio"
import "reflect"

func TestIpacketSerialization(t *testing.T) {
	var src iPacketWire
	var dest iPacketWire

	buffer := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buffer)
	r := bufio.NewReader(buffer)

	src.AtomID = 100
	src.Port.Name = "test_port"
	src.Port.Index = 1
	src.Data = []byte{4, 3, 2, 1}

	err := src.WriteTo(w)
	if err != nil {
		t.Fatal("Serialization failed with error", err)
		return
	}

	dest.ReadFrom(r)
	if err != nil {
		t.Fatal("Deserialization failed with error", err)
		return
	}

	if !reflect.DeepEqual(src, dest) {
		t.Fatal("Source and destination packets are different")
		return
	}
}
