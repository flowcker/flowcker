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
	"sync"

	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("flowcker_atom")

// Atom runs an atom based on an element function passed to it
//
// Usually the process is:
//  + Obtain a IPInbound channel from where all incomming IPs are passed
//     This can be done with a PortsServer and an InboundUnirouter
//  + Create Atom with NewAtom, passing the IPInbound channel and the element
//  + Add inbound middlewares if any are necessary, Standard adds the standard ones
//  + Call Run
//  + Call Setup with the corresponding Atom configuration, this is already done by
//      standard middleware when receiving control IP
//  + Call StartElement, this is done already by standard middleware
//      when receiving control IP
//  + Wait until CloseChannel() channel is closed
type Atom struct {
	started    bool
	configured bool

	inboundPipelineTail chan fc.IPInbound

	toElement    chan fc.IPInbound
	toDispatcher chan fc.IPOutbound
	element      fc.Element

	atom *fc.Atom

	startOnce     sync.Once
	configureOnce sync.Once

	closeChannel chan struct{}
}

// NewAtom creates an atom
//
// Params:
//  + element is the element function of this Atom
//  + input is a channel where the atom is going to
//     receive the IPs
func NewAtom(element fc.Element, input chan fc.IPInbound) (r *Atom) {
	r = new(Atom)

	r.element = element
	r.configured = false
	r.started = false

	r.inboundPipelineTail = input

	r.toElement = make(chan fc.IPInbound)
	r.toDispatcher = make(chan fc.IPOutbound)

	r.closeChannel = make(chan struct{})

	return r
}

// Standard adds standard middleware to Atom
func (r *Atom) Standard() {
	r.UseInbound(AtomSetupMiddleware)
	r.UseInbound(AtomStartMiddleware)
}

// Run starts the loop reading incomming IPs
//
// This should be called after all inbound middleware have been added
func (r *Atom) Run() {
	r.elementIncomingLoop()
}

// UseInbound adds an inbound middleware
func (r *Atom) UseInbound(middleware InboundMiddleware) {
	r.inboundPipelineTail = middleware(r, r.inboundPipelineTail)
}

// GetOutput returns the output channel of the Atom
//
// Usually this is connected to a packet dispatcher
func (r *Atom) GetOutput() chan fc.IPOutbound {
	return r.toDispatcher
}

func (r *Atom) elementIncomingLoop() {
	log.Debug("Starting incomming loop for Atom")
	go func() {
		for incoming := range r.inboundPipelineTail {
			log.Debug("Atom: incoming packet")
			// TODO handling the closing of this
			if !r.started && !r.configured {
				log.Error("Atom received message for non-configured element; missing middleware?")
				spew.Dump(incoming)
				panic("Atom received message for non-configured element; missing middleware?")
			}
			if !r.started && r.configured {
				r.StartElement()
			}
			r.toElement <- incoming
		}
		log.Debug("Mindleware inbound channel closed")
		close(r.toElement)
	}()
}

func (r *Atom) outgoingLoop(input chan fc.IPOutbound) {
	go func() {
		for outgoing := range input {
			log.Debug("Atom: received outbound packet, sending to packet dispatcher...")
			r.toDispatcher <- outgoing
		}
		log.Debug("Atom: Element closed output channel")
		// TODO stops everything if this is closed, perhaps sending control package
		// This signals the end
		close(r.closeChannel)
	}()
}

// StartElement starts the element of this Atom
//
// Setup should be called before calling this function
func (r *Atom) StartElement() {
	log.Debug("Starting element in Atom")
	if !r.configured {
		log.Fatal("Starting element in not configured atom")
		return
	}

	r.startOnce.Do(func() {
		fromElement, _ := r.element(r.atom, r.toElement)
		// TODO, check error
		r.outgoingLoop(fromElement)
		r.started = true
		log.Debug("Atom: element started")
	})
}

// Setup sets up the Atom from a Atom configuration
func (r *Atom) Setup(atom *fc.Atom) {
	log.Debug("Atom seting up")
	r.configureOnce.Do(func() {
		r.atom = atom

		r.configured = true
		log.Debug("Atom setup completed")
	})
}

// CloseChannel returns a channel that closes when the Atom stops
func (r *Atom) CloseChannel() chan struct{} {
	return r.closeChannel
}
