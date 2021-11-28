package messages

import (
	"encoding/binary"
	"io"
)

type Time struct {
	LatencySec  int32
	LatencyUsec int32
}

func CreateTime() Time {
	//now := time.Now()
	second := int32(0)
	usec := int32(0) - second*1000
	return Time{second, usec}

}
func (m Time) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.LittleEndian, m)
	return int64(8), err
}

func (m *Time) ReadFrom(r io.Reader) (int64, error) {
	err := binary.Read(r, binary.LittleEndian, m)
	return int64(8), err
}
