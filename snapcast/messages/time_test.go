package messages

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestTimeWriteTo(t *testing.T) {
	timeValues := make([]int32, 2)
	timeValues[0] = rand.Int31()
	timeValues[1] = rand.Int31()
	time := Time{timeValues[0], timeValues[1]}
	var b bytes.Buffer
	var result1, result2 int32

	_, err := time.WriteTo(&b)
	if err != nil {
		t.Fatalf("time.WriteTo() got error: %s", err)
	}

	reader := bytes.NewReader(b.Bytes())
	err = binary.Read(reader, binary.LittleEndian, &result1)
	if err != nil {
		t.Fatalf("Got error: %s", err)
	}
	if time.LatencySec != result1 {
		t.Errorf("time.LatencySec changed, expected: %d, got: %d", time.LatencySec, result1)
	}

	err = binary.Read(reader, binary.LittleEndian, &result2)
	if err != nil {
		t.Fatalf("Got error: %s", err)
	}

	if time.LatencyUsec != result2 {
		t.Errorf("time.LatencyUsec changed, expected: %d, got: %d", time.LatencyUsec, result2)
	}
}

func TestTimeReadFrom(t *testing.T) {
	var b bytes.Buffer
	timeValues := make([]int32, 2)
	timeValues[0] = rand.Int31()
	timeValues[1] = rand.Int31()
	time := Time{}
	binary.Write(&b, binary.LittleEndian, timeValues[0])
	binary.Write(&b, binary.LittleEndian, timeValues[1])
	reader := bytes.NewReader(b.Bytes())

	_, err := time.ReadFrom(reader)
	if err != nil {
		t.Fatalf("time.ReadFrom() got error: %s", err)
	}

	if time.LatencySec != timeValues[0] {
		t.Errorf("time.LatencySec changed, expected: %d, got: %d", timeValues[0], time.LatencySec)
	}

	if time.LatencyUsec != timeValues[1] {
		t.Errorf("time.LatencyUsec changed, expected: %d, got: %d", timeValues[1], time.LatencyUsec)
	}
}
