package snapcast

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type SnapControl struct {
	host           string
	port           string
	id             string
	volume         uint
	ducking        uint
	duckingState   bool
	requestCounter uint64
}

type controlmessage struct {
	Id      string      `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type volumeparam struct {
	Id     string `json:"id"`
	Volume struct {
		Muted   bool `json:"muted"`
		Percent uint `json:"percent"`
	} `json:"volume"`
}

func CreateSnapControl(host, port, id string, volume, ducking uint) (*SnapControl, error) {
	snapControl := SnapControl{host, port, id, volume, ducking, false, 0}
	err := snapControl.updateRemoteVolume()
	if err != nil {
		return &SnapControl{}, err
	}
	return &snapControl, nil
}

func (s *SnapControl) SetVolume(volume uint) error {
	if volume > 100 {
		return errors.New("volume bigger than 100")
	}

	s.volume = volume
	err := s.updateRemoteVolume()
	return err
}

func (s *SnapControl) SetDucking(ducking uint) error {
	if ducking > 100 {
		return errors.New("ducking volume bigger than 100")
	}

	s.ducking = ducking
	err := s.updateRemoteVolume()
	return err
}

func (s SnapControl) GetVolume() uint {
	return s.volume
}

func (s SnapControl) GetDucking() uint {
	return s.ducking
}

func (s *SnapControl) SetDuckingState(state bool) {
	if s.duckingState != state {
		s.duckingState = state
		s.updateRemoteVolume()
	}
}

func (s *SnapControl) updateRemoteVolume() error {
	vol := s.volume
	if s.duckingState {
		vol = s.ducking
	}
	volPart := volumeparam{}
	volPart.Id = s.id
	volPart.Volume.Muted = false
	volPart.Volume.Percent = vol
	err := s.sendMessage("Client.SetVolume", volPart)
	return err
}

func (s *SnapControl) sendMessage(method string, params interface{}) error {
	url := "http://" + s.host + ":" + s.port + "/jsonrpc"
	counter := strconv.FormatUint(s.requestCounter, 10)
	message := controlmessage{counter, "2.0", method, params}
	s.requestCounter += 1
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("unexpected answer")
	}
	return nil
}
