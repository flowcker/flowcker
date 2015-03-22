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

package common

import "encoding/json"

// Atom token describes an atom (FBP process) in a molecule (FBP graph)
type Atom struct {
	ID      uint32
	Element string
	Config  *json.RawMessage
	Addr    string
}

// Port represent an atom port with index
type Port struct {
	Name  string
	Index uint32
}

// LinkSide: each link has one two of this, on source (From) and destination (To)
type LinkSide struct {
	Port   Port
	AtomID uint32
}

// Link represents the connection between one output port to one input port in different atoms
type Link struct {
	ID   uint32
	From LinkSide
	To   LinkSide
}
