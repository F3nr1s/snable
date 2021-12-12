package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"

	"strconv"
)

func setUpOpusData(t *testing.T) ([]byte, [4]int32) {
	t.Helper()

	// Converting "OPUS" to 4 byte integer
	codecName64, _ := strconv.ParseInt("4F505553", 16, 32)
	codeName := int32(codecName64)
	sampleRate := rand.Int31()
	bitDepth := int16(rand.Int31())
	channel := int16(rand.Int31())
	var b bytes.Buffer

	binary.Write(&b, binary.LittleEndian, uint32(12))
	binary.Write(&b, binary.LittleEndian, codeName)
	binary.Write(&b, binary.LittleEndian, sampleRate)
	binary.Write(&b, binary.LittleEndian, bitDepth)
	binary.Write(&b, binary.LittleEndian, channel)

	return b.Bytes(), [4]int32{codeName, sampleRate, int32(bitDepth), int32(channel)}
}

func compareOpus(opusHeader OpusHeader, data []int32, t *testing.T) {
	t.Helper()
	if opusHeader.Name != "OPUS" {
		t.Errorf("Opusheader.Name: Expected: OPUS, Got: %s", opusHeader.Name)
	}

	if data[1] != opusHeader.SampleRate {
		t.Errorf("Opusheader.SampleRate: Expected: %d, Got: %d", data[1], opusHeader.SampleRate)
	}

	if int16(data[2]) != opusHeader.BitDepth {
		t.Errorf("Opusheader.BitDepth: Expected: %d, Got: %d", int16(data[2]), opusHeader.BitDepth)
	}

	if int16(data[3]) != opusHeader.Channels {
		t.Errorf("Opusheader.Channels: Expected: %d, Got: %d", int16(data[3]), opusHeader.Channels)
	}
}

func TestOpusReadFrom(t *testing.T) {
	var opusHeader OpusHeader
	b, data := setUpOpusData(t)
	reader := bytes.NewReader(b)

	_, err := opusHeader.ReadFrom(reader)
	if err != nil {
		t.Fatalf("OpusHeader.ReadFrom gave error: %s", err)
	}

	compareOpus(opusHeader, data[:], t)
}

func TestCodecReadFromKnownCodec(t *testing.T) {
	opusB, opusData := setUpOpusData(t)
	var codeTests = []struct {
		name string
		b    []byte
		data [4]int32
	}{{OPUS, opusB, opusData}}

	for _, tt := range codeTests {
		var b bytes.Buffer
		var codec Codec
		size := uint32(len(tt.name))
		binary.Write(&b, binary.LittleEndian, size)
		binary.Write(&b, binary.LittleEndian, []byte(tt.name))
		data := append(b.Bytes(), tt.b...)
		reader := bytes.NewReader(data)
		n, err := codec.ReadFrom(reader)
		if err != nil {
			fmt.Print(n)
			t.Fatalf("Codec.ReadFrom gave error: %s", err)
		}

		if tt.name != codec.Codec {
			t.Errorf("Codec.Codec wrong: Expected: %s, Got %s", tt.name, codec.Codec)
		}

		switch codec.Payload.(type) {
		case *OpusHeader:
			o := codec.Payload.(*OpusHeader)
			compareOpus(*o, tt.data[:], t)
		default:
			t.Fatalf("Unknown Payload Type %T", codec.Payload)
		}
	}
}

func TestCodeReadFromUnknownCodec(t *testing.T) {
	codecName := "Unknown"
	payload := rand.Int63()
	size := uint32(len(codecName))
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, size)
	binary.Write(&b, binary.LittleEndian, []byte(codecName))
	binary.Write(&b, binary.LittleEndian, uint32(4))
	binary.Write(&b, binary.LittleEndian, payload)
	var codec Codec
	reader := bytes.NewReader(b.Bytes())

	_, err := codec.ReadFrom(reader)
	if err == nil {
		t.Fatal("No error returned")
	}

	if err.Error() != "Unknown/Handled header type" {
		t.Errorf("Wrong Error code, got: %v", err)
	}
}
