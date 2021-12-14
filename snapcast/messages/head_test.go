package messages

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestGetMessageTypeName(t *testing.T) {
	cases := []struct {
		id       uint16
		expected string
	}{
		{0, "Base"},
		{1, "Codec"},
		{2, "WireChunk"},
		{3, "ServerSetting"},
		{4, "Time"},
		{5, "Hello"},
		{6, "StreamTag"},
	}

	for _, tt := range cases {
		got := GetMessageTypeName(tt.id)
		if got != tt.expected {
			t.Errorf("Got wrong Type, Expected: %s, Got: %s", tt.expected, got)
		}
	}
}

func TestHeadReadFrom(t *testing.T) {
	msgType := uint16(rand.Intn(7))
	id := uint16(rand.Int())
	refersTo := uint16(rand.Int())
	sent_sec := rand.Int31()
	sent_usec := rand.Int31()
	received_sec := rand.Int31()
	received_usec := rand.Int31()
	size := uint32(rand.Int31())
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, msgType)
	binary.Write(&b, binary.LittleEndian, id)
	binary.Write(&b, binary.LittleEndian, refersTo)
	binary.Write(&b, binary.LittleEndian, sent_sec)
	binary.Write(&b, binary.LittleEndian, sent_usec)
	binary.Write(&b, binary.LittleEndian, received_sec)
	binary.Write(&b, binary.LittleEndian, received_usec)
	binary.Write(&b, binary.LittleEndian, size)
	reader := bytes.NewReader(b.Bytes())
	var head Head

	_, err := head.ReadFrom(reader)
	if err != nil {
		t.Fatalf("Head.ReadFrom() gave error: %s", err)
	}

	if msgType != head.MsgType {
		t.Errorf("Wrong MsgType, Expected: %d, Got: %d", msgType, head.MsgType)
	}

	if id != head.Id {
		t.Errorf("Wrong Id, Expected: %d, Got: %d", id, head.Id)
	}

	if refersTo != head.RefersTo {
		t.Errorf("Wrong RefersTo, Expected: %d, Got: %d", refersTo, head.RefersTo)
	}

	if sent_sec != head.Sent_sec {
		t.Errorf("Wrong Sent_sec, Expected: %d, Got: %d", sent_sec, head.Sent_sec)
	}

	if sent_usec != head.Sent_usec {
		t.Errorf("Wrong Sent_usec, Expected: %d, Got: %d", sent_usec, head.Sent_usec)
	}

	if received_sec != head.Received_sec {
		t.Errorf("Wrong Received_sec, Expected: %d, Got: %d", received_sec, head.Received_sec)
	}

	if received_usec != head.Received_usec {
		t.Errorf("Wrong Received_usec, Expected: %d, Got: %d", received_usec, head.Received_usec)
	}

	if size != head.Size {
		t.Errorf("Wrong Size, Expected: %d, Got: %d", size, head.Size)
	}
}

func TestHeadWriteTo(t *testing.T) {
	msgType := uint16(rand.Intn(7))
	id := uint16(rand.Int())
	refersTo := uint16(rand.Int())
	sent_sec := rand.Int31()
	sent_usec := rand.Int31()
	received_sec := rand.Int31()
	received_usec := rand.Int31()
	size := uint32(rand.Int31())
	var msgResult, idResult, refersToResult uint16
	var sent_secResult, sent_usecResult, received_secResult, received_usecResult int32
	var sizeResult uint32
	var b bytes.Buffer
	head := Head{
		MsgType:       msgType,
		Id:            id,
		RefersTo:      refersTo,
		Sent_sec:      sent_sec,
		Sent_usec:     sent_usec,
		Received_sec:  received_sec,
		Received_usec: received_usec,
		Size:          size,
	}

	_, err := head.WriteTo(&b)
	if err != nil {
		t.Fatalf("Head.WriteTo() gave error: %s", err)
	}

	reader := bytes.NewReader(b.Bytes())
	binary.Read(reader, binary.LittleEndian, &msgResult)
	if msgResult != msgType {
		t.Errorf("MsgType changed: Expected: %d, Got: %d", msgType, msgResult)
	}

	binary.Read(reader, binary.LittleEndian, &idResult)
	if idResult != id {
		t.Errorf("Id changed: Expected: %d, Got: %d", id, idResult)
	}

	binary.Read(reader, binary.LittleEndian, &refersToResult)
	if refersToResult != refersTo {
		t.Errorf("RefersTo changed: Expected: %d, Got: %d", refersTo, refersToResult)
	}

	binary.Read(reader, binary.LittleEndian, &sent_secResult)
	if sent_secResult != sent_sec {
		t.Errorf("Sent_sec changed: Expected: %d, Got: %d", sent_sec, sent_secResult)
	}

	binary.Read(reader, binary.LittleEndian, &sent_usecResult)
	if sent_usecResult != sent_usec {
		t.Errorf("Sent_usec changed: Expected: %d, Got: %d", sent_usec, sent_usecResult)
	}

	binary.Read(reader, binary.LittleEndian, &received_secResult)
	if received_secResult != received_sec {
		t.Errorf("Received_sec changed: Expected: %d, Got: %d", received_sec, received_secResult)
	}

	binary.Read(reader, binary.LittleEndian, &received_usecResult)
	if received_usecResult != received_usec {
		t.Errorf("Received_usec changed: Expected: %d, Got: %d", received_usec, received_usecResult)
	}

	binary.Read(reader, binary.LittleEndian, &sizeResult)
	if sizeResult != size {
		t.Errorf("Size changed: Expected: %d, Got: %d", size, sizeResult)
	}
}

func TestHeadReadFromFailure(t *testing.T) {
	var b []byte
	r := bytes.NewReader(b)
	var h Head

	n, err := h.ReadFrom(r)
	if err == nil {
		t.Errorf("Got no err, read %d bytes", n)
	}
	if err.Error() != "EOF" {
		t.Errorf("Unexpected error: %s", err)
	}
}
