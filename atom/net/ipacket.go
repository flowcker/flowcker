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
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
)

type iPacketWire struct {
	fc.IPacket
	fc.LinkSide
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func readZstring(reader *bufio.Reader) (res string, err error) {
	var raw []byte

	raw, err = reader.ReadBytes(0)
	if err != nil {
		return "", err
	}

	if raw[len(raw)-1] != 0 {
		panic("iPacketWire: not zstring found in readZstring")
	}

	raw = raw[:len(raw)-1]
	return string(raw), nil
}

func readUint32(reader *bufio.Reader) (val uint32, err error) {
	var raw [4]byte
	var n int

	n, err = io.ReadFull(reader, raw[:])
	if n != 4 {
		panic("iPacketWire: wrong number of bytes read in readUint32")
	}
	if err != nil {
		return 0, err
	}

	buf := bytes.NewBuffer(raw[:])
	err = binary.Read(buf, binary.BigEndian, &val)

	return val, err
}

func writeZstring(w *bufio.Writer, s string) (err error) {
	n, err := io.WriteString(w, s)
	if n < len(s) {
		panic("not all line written")
	}
	checkError(err)
	err = w.WriteByte(0)
	checkError(err)

	return err
}

func writeUint32(w *bufio.Writer, val uint32) (err error) {
	err = binary.Write(w, binary.BigEndian, val)

	return err
}

// ReadFrom reads from a bufio an information packet
func (p *iPacketWire) ReadFrom(reader *bufio.Reader) (err error) {
	var header [4]byte

	_, err = io.ReadFull(reader, header[:])
	if err != nil {
		return err
	}

	if string(header[:]) != "FCIP" {
		log.Error("Incorrect IP header \n%s", spew.Sdump(header))
		panic("Incorrect IP header")
	}

	pType, err := readUint32(reader)
	if err != nil {
		return err
	}
	p.Type = fc.IPType(pType)

	if p.AtomID, err = readUint32(reader); err != nil {
		return err
	}

	if p.Port.Name, err = readZstring(reader); err != nil {
		return err
	}

	if p.Port.Index, err = readUint32(reader); err != nil {
		return err
	}

	var dataLength uint32
	if dataLength, err = readUint32(reader); err != nil {
		return err
	}

	if dataLength != 0 {
		p.Data = make([]byte, dataLength)
		n, err := io.ReadFull(reader, p.Data)
		if n != int(dataLength) {
			log.Warning("IPacketWire: read less payload than required")
			return errors.New("IPacketWire: read less payload than required")
		}
		if err != nil {
			log.Error("IPacketWire: error reading payload")
			return err
		}
	}

	var tailing [4]byte

	_, err = io.ReadFull(reader, tailing[:])
	if err != nil {
		return err
	}

	if string(tailing[:]) != "PICF" {
		log.Error("Incorrect IP tailing \n%s\n%s\nDataLength: %d\n", spew.Sdump(tailing), spew.Sdump(p), dataLength)
		panic("Incorrect IP tailing")
	}

	return nil
}

// WriteTo writes to a bufio an information packet
func (p *iPacketWire) WriteTo(w *bufio.Writer) (err error) {
	_, err = w.Write([]byte("FCIP"))
	checkError(err)

	err = writeUint32(w, uint32(p.Type))
	checkError(err)

	err = writeUint32(w, p.AtomID)
	checkError(err)

	err = writeZstring(w, p.Port.Name)
	checkError(err)

	err = writeUint32(w, p.Port.Index)
	checkError(err)

	err = writeUint32(w, uint32(len(p.Data)))
	checkError(err)

	if len(p.Data) != 0 {
		var n int
		n, err = w.Write(p.Data)
		if n != len(p.Data) {
			panic("iPacketWire: Not writen all data")
		}
		checkError(err)
	}

	_, err = w.Write([]byte("PICF"))
	checkError(err)

	w.Flush()

	return
}

func newIPWire(data []byte, port fc.Port, atomID uint32) *iPacketWire {
	return &iPacketWire{fc.IPacket{Type: fc.DATA, Data: data}, fc.LinkSide{AtomID: atomID, Port: port}}
}

func newIPWireType(t fc.IPType, data []byte, port fc.Port, atomID uint32) *iPacketWire {
	return &iPacketWire{fc.IPacket{Type: t, Data: data}, fc.LinkSide{AtomID: atomID, Port: port}}
}
