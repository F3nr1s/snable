package messages

import (
	"encoding/binary"
	"io"
)

type WireChunk struct {
	TimestampSec  int32
	TimestampUsec int32
	Payload       string
}

func (m *WireChunk) ReadFrom(r io.Reader) (int64, error) {
	var timestamp int32
	var size uint32
	var n int64 = 0
	err := binary.Read(r, binary.LittleEndian, &timestamp)
	if err != nil {
		return n, err
	}
	n += 4
	m.TimestampSec = timestamp
	err = binary.Read(r, binary.LittleEndian, &timestamp)
	if err != nil {
		return n, err
	}
	n += 4
	m.TimestampUsec = timestamp
	err = binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return n, err
	}
	n += 4
	b := make([]byte, size)
	err = binary.Read(r, binary.LittleEndian, b)
	if err != nil {
		return n, err
	}
	m.Payload = string(b)
	return n + int64(size), err
}
