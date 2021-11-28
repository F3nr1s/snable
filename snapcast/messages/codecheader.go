package messages

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	OPUS = "opus"
)

type Header interface {
	ReadFrom(io.Reader) (int64, error)
}

type OpusHeader struct {
	Name       string
	SampleRate int32
	BitDepth   int16
	Channels   int16
}

func (h *OpusHeader) ReadFrom(r io.Reader) (int64, error) {
	var n int64
	var size uint32

	// Size
	err := binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return n, err
	}
	n += 4

	// Name
	b := make([]byte, 4)
	err = binary.Read(r, binary.LittleEndian, b)
	if err != nil {
		return n, err
	}
	n += 4
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-1-i] = b[len(b)-1-i], b[i]
	}

	h.Name = string(b)

	// Sample rate
	err = binary.Read(r, binary.LittleEndian, &h.SampleRate)
	if err != nil {
		return n, err
	}
	n += 4

	err = binary.Read(r, binary.LittleEndian, &h.BitDepth)
	if err != nil {
		return n, err
	}
	n += 2

	//Channels
	err = binary.Read(r, binary.LittleEndian, &h.Channels)
	if err != nil {
		return n, err
	}
	n += 2

	return n, nil
}

type Codec struct {
	Codec   string
	Payload Header
}

func (m *Codec) ReadFrom(r io.Reader) (int64, error) {
	var size uint32
	var n int64 = 0
	//var payload Header
	err := binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return n, err
	}
	n += 4
	b := make([]byte, size)
	err = binary.Read(r, binary.LittleEndian, b)
	if err != nil {
		return n, err
	}
	n += int64(len(b))
	m.Codec = string(b)
	switch m.Codec {
	case OPUS:
		m.Payload = new(OpusHeader)
		a, err := m.Payload.ReadFrom(r)
		if err != nil {
			return n, err
		}
		n += a
	default:
		return n, errors.New("Unknown/Handled header type")
	}

	return n, err
}
