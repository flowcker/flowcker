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

const (
	// DATA - IP contains data
	DATA IPType = 0
	// DISCONNECT - Link is closed and source will not send anymore
	DISCONNECT IPType = 13
	// CONTROL - IP contains control data
	CONTROL IPType = 20
)

// IPType contains IP type
type IPType uint32

// IP interface for information packages
type IP interface {
	GetType() IPType
	GetData() []byte
}

// IPInbound is a IP that is incoming to the Atom
type IPInbound interface {
	IP
	GetTo() Port
}

// IPOutbound is a IP that is outgoing from the Atom
type IPOutbound interface {
	IP
	GetFrom() Port
	GetAll() bool
	GetIndexSelected() bool
}
