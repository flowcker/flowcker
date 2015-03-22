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
	"encoding/json"
	"io/ioutil"
)

// Molecule represents a flowcker program (FBP graph)
type Molecule struct {
	Atoms []Atom
	Links []Link
}

// NewMoleculeFromFile creates a molecule from a file
func NewMoleculeFromFile(filename string) *Molecule {
	mol := new(Molecule)

	blob, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err) // TODO
	}

	mol.importJSON(blob)
	mol.indexLinks()

	return mol
}

func (mol *Molecule) importJSON(blob []byte) (err error) {
	err = json.Unmarshal(blob, &mol)
	return
}

// AddAtom adds an atom to the molecule
func (mol *Molecule) AddAtom(element string, config string, addr string) {
	ID := mol.maxAtomID()

	if config == "" {
		config = "{}"
	}

	var configraw json.RawMessage
	configraw.UnmarshalJSON([]byte(config))

	mol.Atoms = append(mol.Atoms, Atom{ID: ID, Element: element, Config: &configraw, Addr: addr})
}

type lookupT map[uint32]map[string]uint32

func (mol *Molecule) indexLinks() {
	outportsC := make(lookupT)
	inportsC := make(lookupT)

	for k, link := range mol.Links {
		mol.Links[k].From.Port.Index = indexLinksLookup(outportsC, link.From.AtomID, link.From.Port.Name)
		mol.Links[k].To.Port.Index = indexLinksLookup(inportsC, link.To.AtomID, link.To.Port.Name)
	}
}

func indexLinksLookup(t lookupT, atomID uint32, name string) uint32 {
	l1, ok1 := t[atomID]
	if !ok1 {
		t[atomID] = make(map[string]uint32)
	}
	l2, ok2 := l1[name]
	if !ok2 {
		t[atomID][name] = 0
		return 0
	}

	t[atomID][name] = l2 + 1
	return t[atomID][name]
}

func (mol *Molecule) maxAtomID() uint32 {
	var max uint32 // Note that it is initalized to zero by Go
	for _, atom := range mol.Atoms {
		if atom.ID > max {
			max = atom.ID
		}
	}

	return max
}
