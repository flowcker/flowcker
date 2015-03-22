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

package atom

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/flowcker/flowcker/atom/net"
	fc "github.com/flowcker/flowcker/common"
)

// InboundMiddleware interface for any middleware that can be used in the inbound pipeline of an Atom
type InboundMiddleware func(*Atom, chan fc.IPInbound) chan fc.IPInbound

func AtomSetupMiddleware(r *Atom, input chan fc.IPInbound) chan fc.IPInbound {
	output := make(chan fc.IPInbound)
	findAtom := func(atomID uint32, atoms []fc.Atom) (result *fc.Atom) {
		for _, atom := range atoms {
			if atom.ID == atomID {
				return &atom
			}
		}

		return nil
	}

	go func() {
		for incoming := range input {
			if incoming.GetTo().Name == "control" {
				log.Debug("AtomSetupMiddleware: received control message")
				var msg ControlMessage

				err := json.Unmarshal(incoming.GetData(), &msg)
				if err != nil {
					spew.Dump(incoming)
					panic("Error unmarshalling control package")
				}
				if msg.Type == "AtomSetup" {
					log.Debug("AtomSetupMiddleware: received AtomSetup")
					var atomSetup AtomSetup
					err := json.Unmarshal(msg.Data, &atomSetup)
					if err != nil {
						spew.Dump(msg.Data)
						panic("Error unmarshalling AtomSetup package")
					}
					atomDispatcherConf := net.NewPacketDispatcherConf(atomSetup.AtomID, atomSetup.Atoms, atomSetup.Links)

					atom := findAtom(atomSetup.AtomID, atomSetup.Atoms)
					net.NewPacketDispatcher(atomDispatcherConf, r.GetOutput())
					r.Setup(atom)

					continue // Do not pass it through in this case
				}
			}
			output <- incoming
		}
		close(output)
	}()

	return output
}

func AtomStartMiddleware(r *Atom, input chan fc.IPInbound) chan fc.IPInbound {
	output := make(chan fc.IPInbound)

	go func() {
		for incoming := range input {
			if incoming.GetTo().Name == "control" {
				log.Debug("AtomStartMiddleware: received control message")
				var msg ControlMessage

				msg.JSONDecode(incoming.GetData())
				if msg.Type == "AtomStart" {
					log.Debug("AtomStartMiddleware: received AtomStart")
					r.StartElement()
					continue // Do not pass it through in this case
				}
			}
			output <- incoming
		}
		close(output)
	}()

	return output
}
