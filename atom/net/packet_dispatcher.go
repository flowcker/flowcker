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
	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
)

// PacketDispatcher allows an Atom to send information packages to other atoms via PortClient
type PacketDispatcher struct {
	conf            *PacketDispatcherConf
	atomsD          map[uint32]*atomDispatcher
	portsD          map[string]map[uint32]*portDispatcher
	portsDinterator map[string]chan uint32

	input chan fc.IPOutbound
}

// NewPacketDispatcher creates a packet from a AtomDispatcherConf
func NewPacketDispatcher(conf *PacketDispatcherConf, input chan fc.IPOutbound) (pd *PacketDispatcher) {
	pd = new(PacketDispatcher)
	pd.conf = conf
	pd.input = input
	pd.createDispatchers()
	pd.start()

	return pd
}

// Dispatch a package, routes it to the right atom and destiantion port
func (pd *PacketDispatcher) Dispatch(p fc.IPOutbound) {
	outPortName := p.GetFrom().Name
	if pd.portsD[outPortName] == nil {
		log.Debug("Dispatching packet to unkown port %s, discarded", outPortName)
		return
	}

	if p.GetAll() {
		pd.dispatchAll(p)
		return
	}

	if !p.GetIndexSelected() {
		p = pd.dispatchIterator(p)
	}
	outPortIndex := p.GetFrom().Index

	log.Debug("PacketDispatcher: dispatching packet to %s[%d] with data \n%s\n%s", outPortName, outPortIndex, spew.Sdump(p.GetData()), spew.Sdump(p))

	err := pd.portsD[outPortName][outPortIndex].dispatch(p)
	if err != nil {
		log.Error("Error dispatching IP")
		panic(err)
	}
}

func (pd *PacketDispatcher) dispatchAll(p fc.IPOutbound) {
	log.Debug("PacketDispatcher: dispatching packet to %s[All]", p.GetFrom().Name)
	for index := range pd.portsD[p.GetFrom().Name] {
		pd.Dispatch(fc.NewIPOut(p, fc.Port{Name: p.GetFrom().Name, Index: index}))
	}
}

func (pd *PacketDispatcher) dispatchIterator(p fc.IPOutbound) fc.IPOutbound {
	outPortName := p.GetFrom().Name
	if pd.portsDinterator[outPortName] == nil {
		log.Debug("creating iterator for port %s", outPortName)
		// TODO with an slice of the keys
		iter := make(chan uint32, 10)
		go func() {
			for {
				for k := range pd.portsD[outPortName] {
					log.Debug("generated index %d, for port %s", k, outPortName)
					iter <- k
				}
			}
		}()
		pd.portsDinterator[outPortName] = iter
	}

	log.Debug("PacketDispatcher: selecting index for packet to %s, \n%s", outPortName, spew.Sdump(p))
	index := <-pd.portsDinterator[outPortName]

	return fc.NewIPOut(p, fc.Port{Name: outPortName, Index: index})
}

func (pd *PacketDispatcher) createDispatchers() {
	pd.portsDinterator = make(map[string]chan uint32)
	pd.atomsD = make(map[uint32]*atomDispatcher)
	for atomID, atomAddr := range pd.conf.atomsAddr {
		pd.atomsD[atomID] = newAtomDispatcher(atomID, atomAddr)
	}

	pd.portsD = make(map[string]map[uint32]*portDispatcher)
	for outPortName, convs := range pd.conf.portsConv {
		pd.portsD[outPortName] = make(map[uint32]*portDispatcher)
		for outPortIndex, conv := range convs {
			pd.portsD[outPortName][outPortIndex] = newPortDispatcher(conv.AtomID, conv.Port.Name, conv.Port.Index, pd.atomsD[conv.AtomID])
		}
	}
}

func (pd *PacketDispatcher) start() {
	go func() {
		for incoming := range pd.input {
			log.Debug("PacketDispatcher: recevied outbound packet")
			pd.Dispatch(incoming)
		}
	}()
}
