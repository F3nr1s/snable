package messages

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

type ServerSettings struct {
	BufferMs int  `json:"bufferMs"`
	Latency  int  `json:"latency"`
	Muted    bool `json:"muted"`
	Volume   int  `json:"volume"`
}

func (m *ServerSettings) ReadFrom(r io.Reader) (int64, error) {
	var size uint32
	err := binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return 0, nil
	}

	var n int64 = 4

	b := make([]byte, size)
	err = binary.Read(r, binary.LittleEndian, b)
	if err != nil {
		return n, err
	}
	err = json.Unmarshal(b, m)
	return int64(size) + n, err
}
