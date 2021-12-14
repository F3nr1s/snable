package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
)

type Hello struct {
	Arch                      string `json:"Arch"`
	ClientName                string `json:"ClientName"`
	HostName                  string `json:"HostName"`
	ID                        string `json:"ID"`
	Instance                  int    `json:"Instance"`
	Mac                       string `json:"MAC"`
	Os                        string `json:"OS"`
	SnapStreamProtocolVersion int    `json:"SnapStreamProtocolVersion"`
	Version                   string `json:"Version"`
}

func (m Hello) Size() uint32 {
	j, _ := json.Marshal(m)
	return uint32(len(j))
}

func (m Hello) FullSize() uint32 {
	n := m.Size()
	return uint32(4 + n)
}

func (m Hello) WriteTo(w io.Writer) (int64, error) {
	bodySize := m.Size()
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, bodySize)
	msg, _ := json.Marshal(m)

	binary.Write(b, binary.LittleEndian, msg)
	s, err := w.Write(b.Bytes())

	return int64(s), err
}
