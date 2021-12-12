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

func (m Hello) Size() (uint32, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}

	return uint32(len(j)), err
}

func (m Hello) FullSize() (uint32, error) {
	n, err := m.Size()
	if err != nil {
		return 0, err
	}

	return uint32(4 + n), err
}

func (m Hello) WriteTo(w io.Writer) (int64, error) {
	bodySize, err := m.Size()
	if err != nil {
		return 0, err
	}
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, bodySize)
	msg, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}
	binary.Write(b, binary.LittleEndian, msg)
	s, err := w.Write(b.Bytes())

	return int64(s), err
}
