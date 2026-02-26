package protocol

import (
	"bytes"
	"encoding/binary"

	"google.golang.org/protobuf/proto"
)

type Frame struct {
	packets []*Packet
}

func NewFrame() *Frame {
	return &Frame{
		packets: []*Packet{},
	}
}

func BuildFrame(data []byte) (*Frame, error) {
	packets := []*Packet{}
	reader := bytes.NewReader(data)

	for reader.Len() > 0 {
		var length uint16
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		rawPacket := make([]byte, length)
		if _, err := reader.Read(rawPacket); err != nil {
			return nil, err
		}

		var packet *Packet
		if err := proto.Unmarshal(rawPacket, packet); err != nil {
			return nil, err
		}

		packets = append(packets, packet)
	}

	return &Frame{
		packets: packets,
	}, nil
}

func (f *Frame) Append(packet *Packet) {
	f.packets = append(f.packets, packet)
}

func (f *Frame) Packets() []*Packet {
	return f.packets
}

func (f *Frame) Len() int {
	return len(f.packets)
}

func (f *Frame) Empty() {
	f.packets = f.packets[:0]
}

func (f *Frame) Bytes() ([]byte, error) {
	buf := []byte{}

	for _, packet := range f.packets {
		bytes, err := proto.Marshal(packet)
		if err != nil {
			return nil, err
		}

		length := uint16(len(bytes))
		buf := make([]byte, length)

		binary.BigEndian.PutUint16(buf[:2], length)
		copy(buf[2:], bytes)

		buf = append(buf, buf...)
	}

	return buf, nil
}
