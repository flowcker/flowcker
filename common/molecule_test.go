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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicMoleculeJsonImport(t *testing.T) {
	mol := NewMoleculeFromFile("./test.fm")

	assert.Equal(t, len(mol.Atoms), 4)
	assert.Equal(t, len(mol.Links), 3)
}
