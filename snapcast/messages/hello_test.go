package messages

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestHelloWriteTo(t *testing.T) {
	message := Hello{
		Arch:                      "12345",
		ClientName:                "Testing",
		ID:                        "6789",
		Instance:                  1,
		Mac:                       "c0:ff:ee:c0:ff:ee",
		Os:                        "Nothing",
		SnapStreamProtocolVersion: 3,
		Version:                   "3.1",
	}

	var b bytes.Buffer
	s, err := message.WriteTo(&b)
	if err != nil {
		t.Fatalf("Hello.WriteTo() gave error: %s", err)
	}
	sizeWanted, err := message.FullSize()
	if err != nil {
		t.Fatalf("Hello.Fullsize() gave error: %s", err)
	}
	if s != int64(sizeWanted) {
		t.Errorf("Size is incorrect, expected %d, got %d", sizeWanted, s)
	}

	got := make([]byte, sizeWanted-4)
	msg, _ := json.Marshal(message)
	reader := bytes.NewReader(b.Bytes())
	_, err = reader.ReadAt(got, 4)
	if err != nil {
		t.Fatalf("Reading gave error: %s", err)
	}
	if string(msg) != string(got) {
		t.Errorf("Message changed: expected %s, got %s", msg, got)
	}
}
