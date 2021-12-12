package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
)

func TestServerSettingsReadFrom(t *testing.T) {
	bufferMs := rand.Intn(2500)
	latency := rand.Intn(2500)
	volume := rand.Intn(101)
	data := fmt.Sprintf("{\"bufferMs\": %d, \"latency\": %d, \"muted\": false, \"volume\": %d}", bufferMs, latency, volume)
	size := int32(len(data))
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, size)
	binary.Write(&b, binary.LittleEndian, []byte(data))
	reader := bytes.NewReader(b.Bytes())
	var ServerSettingMsg ServerSettings

	n, err := ServerSettingMsg.ReadFrom(reader)
	if err != nil {
		t.Fatalf("ServerSettings.ReadFrom() gave error: %s, with %d bytes", err, n)
	}
	if n != int64(size+4) {
		t.Errorf("Read wrong amount, expected: %d, got: %d", size+4, n)
	}
	if bufferMs != ServerSettingMsg.BufferMs {
		t.Errorf("ServerSettings.BufferMs changed: expected: %d, got %d", bufferMs, ServerSettingMsg.BufferMs)
	}

	if latency != ServerSettingMsg.Latency {
		t.Errorf("ServerSettings.Latency changed: expected: %d, got %d", latency, ServerSettingMsg.Latency)
	}

	if false != ServerSettingMsg.Muted {
		t.Errorf("ServerSettings.muted changed: expected: %v, got %v", false, ServerSettingMsg.Muted)
	}

	if volume != ServerSettingMsg.Volume {
		t.Errorf("ServerSettings.Volume changed: expected: %d, got %d", volume, ServerSettingMsg.Volume)
	}
}
