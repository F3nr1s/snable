package messages

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	BaseMsg = iota
	CodecMsg
	WireChunkMsg
	ServerSettingMsg
	TimeMsg
	HelloMsg
	StreamTagMsg
)

func GetMessageTypeName(i uint16) string {
	return [...]string{"Base", "Codec", "WireChunk", "ServerSetting", "Time", "Hello", "StreamTag"}[i]
}

type message interface {
	io.WriterTo
	io.ReaderFrom
}

type Head struct {
	MsgType       uint16
	Id            uint16
	RefersTo      uint16
	Sent_sec      int32
	Sent_usec     int32
	Received_sec  int32
	Received_usec int32
	Size          uint32
}

func (m Head) WriteTo(w io.Writer) (int64, error) {
	//now := time.Now()
	b := new(bytes.Buffer)
	/*if m.Sent_sec == 0 {
		m.Sent_sec = int32(now.UTC().Unix())
	}

	if m.Sent_usec == 0 {
		m.Sent_usec = int32(now.UTC().UnixMicro() % 1000)
	}*/
	err := binary.Write(b, binary.LittleEndian, m)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(b.Bytes())
	return int64(n), err
}

func (m *Head) ReadFrom(r io.Reader) (int64, error) {
	err := binary.Read(r, binary.LittleEndian, m)
	if err != nil {
		return 0, err
	}
	return int64(26), err
}
