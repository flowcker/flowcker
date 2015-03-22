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
	"errors"

	fc "github.com/flowcker/flowcker/common"
)

type PacketDispatcherConf struct {
	portsConv map[string]map[uint32]fc.LinkSide
	atomsAddr map[uint32]string
}

func NewPacketDispatcherConf(atomID uint32, atoms []fc.Atom, links []fc.Link) (conf *PacketDispatcherConf) {
	conf = new(PacketDispatcherConf)
	conf.portsConv = make(map[string]map[uint32]fc.LinkSide)
	conf.atomsAddr = make(map[uint32]string)
	for _, link := range links {
		if link.From.AtomID != atomID {
			continue
		}

		if _, already := conf.portsConv[link.From.Port.Name]; !already {
			conf.portsConv[link.From.Port.Name] = make(map[uint32]fc.LinkSide)
		}
		conf.portsConv[link.From.Port.Name][link.From.Port.Index] = link.To
		conf.addAtomAddr(link.To.AtomID, atoms)
	}

	return conf
}

func (conf *PacketDispatcherConf) addAtomAddr(atomID uint32, atoms []fc.Atom) (err error) {
	if _, ok := conf.atomsAddr[atomID]; ok {
		return
	}

	var addr string

	for _, atom := range atoms {
		if atom.ID == atomID {
			addr = atom.Addr
			break
		}
	}

	if addr == "" {
		return errors.New("Atom not found")
	}

	conf.atomsAddr[atomID] = addr

	return nil
}
