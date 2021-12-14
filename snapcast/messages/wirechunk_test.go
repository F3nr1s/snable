package messages

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestWireChunkReadFrom(t *testing.T) {
	var b bytes.Buffer
	var message WireChunk
	size := uint32(rand.Intn(256))
	token := make([]byte, size)
	timeValue := make([]int32, 2)
	timeValue[0] = rand.Int31()
	timeValue[1] = rand.Int31()
	rand.Read(token)
	binary.Write(&b, binary.LittleEndian, timeValue[0])
	binary.Write(&b, binary.LittleEndian, timeValue[1])
	binary.Write(&b, binary.LittleEndian, size)
	binary.Write(&b, binary.LittleEndian, token)

	reader := bytes.NewReader(b.Bytes())
	n, err := message.ReadFrom(reader)

	if err != nil {
		t.Fatalf("WireChunk.ReadFrom() gave error: %s, after %d bytes", err, n)
	}
	if message.TimestampSec != timeValue[0] {
		t.Errorf("Wrong TimeStampSec, expected: %d, got: %d", timeValue[0], message.TimestampSec)
	}
	if message.TimestampUsec != timeValue[1] {
		t.Errorf("Wrong TimeStampUsec, expected: %d, got: %d", timeValue[1], message.TimestampUsec)
	}
	if message.Payload != string(token) {
		t.Error("Wrong Payload")
	}
}

func TestWireChunkFailures(t *testing.T) {
	var b bytes.Buffer
	var message WireChunk
	reader := bytes.NewReader(b.Bytes())
	t.Run("EOF while reading Timestampsec", func(t *testing.T) {
		n, err := message.ReadFrom(reader)
		if err == nil {
			t.Errorf("Expected no read, read %d bytes", n)
		}
		if err.Error() != "EOF" {
			t.Errorf("Unexpected Error: %v", err)
		}
	})

	time32 := rand.Int31()
	binary.Write(&b, binary.LittleEndian, time32)
	reader = bytes.NewReader(b.Bytes())
	t.Run("EOF while reading Timestampusec", func(t *testing.T) {
		n, err := message.ReadFrom(reader)
		if err == nil {
			t.Errorf("Expected no read, read %d bytes", n)
		}
		if err.Error() != "EOF" {
			t.Errorf("Unexpected Error: %v", err)
		}
		if n != 4 {
			t.Errorf("Not right amount of bytes read. Expected: 4, Got: %d", n)
		}
	})

	binary.Write(&b, binary.LittleEndian, time32)
	reader = bytes.NewReader(b.Bytes())
	t.Run("EOF while reading Size", func(t *testing.T) {
		n, err := message.ReadFrom(reader)
		if err == nil {
			t.Errorf("Expected no read, read %d bytes", n)
		}
		if err.Error() != "EOF" {
			t.Errorf("Unexpected Error: %v", err)
		}
		if n != 8 {
			t.Errorf("Not right amount of bytes read. Expected: 8, Got: %d", n)
		}
	})

	size := uint32(rand.Intn(256))
	binary.Write(&b, binary.LittleEndian, size)
	reader = bytes.NewReader(b.Bytes())
	t.Run("EOF while reading Payload", func(t *testing.T) {
		n, err := message.ReadFrom(reader)
		if err == nil {
			t.Errorf("Expected no read, read %d bytes", n)
		}
		if err.Error() != "EOF" {
			t.Errorf("Unexpected Error: %v", err)
		}
		if n != 12 {
			t.Errorf("Not right amount of bytes read. Expected: 12, Got: %d", n)
		}
	})
	token := make([]byte, size-1)
	rand.Read(token)
	binary.Write(&b, binary.LittleEndian, token)
	reader = bytes.NewReader(b.Bytes())
	t.Run("EOF while reading to small payload", func(t *testing.T) {
		n, err := message.ReadFrom(reader)
		if err == nil {
			t.Errorf("Expected no read, read %d bytes", n)
		}
		if err.Error() != "unexpected EOF" {
			t.Errorf("Unexpected Error: %v", err)
		}
		if n != 12 {
			t.Errorf("Not right amount of bytes read. Expected: 12, Got: %d", n)
		}
	})
}
