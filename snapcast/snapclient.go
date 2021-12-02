package snapcast

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"

	"github.com/f3nr1s/snable/snapcast/messages"

	"github.com/matishsiao/goInfo"
	"github.com/sirupsen/logrus"
	"gopkg.in/hraban/opus.v2"
)

type decoder func([]byte, messages.ServerSettings) ([]int16, error)

type SnapClient struct {
	conn          net.Conn
	id            string
	log           *logrus.Logger
	serverSetting messages.ServerSettings
	codecHeader   messages.Codec
	timeMessage   messages.Time
	timeHead      messages.Head
	latencySec    int32
	latencyUsec   int32
	initialized   bool
	output        chan []int16
	d             decoder
}

func Create(host, port, id string, output chan []int16, log *logrus.Logger) (SnapClient, error) {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		//logErrWriter.Println(err)
		return SnapClient{}, err
	}
	client := SnapClient{
		conn,
		id,
		log,
		messages.ServerSettings{},
		messages.Codec{},
		messages.Time{},
		messages.Head{},
		0,
		0,
		false,
		output,
		nil}
	return client, err
}

func (client SnapClient) Close() {
	client.conn.Close()
}

func (client *SnapClient) Initialize() error {
	gi, _ := goInfo.GetInfo()
	msg := messages.Hello{"x86_64", "snable", gi.Hostname, client.id, 1, client.id, gi.GoOS, 2, "0.17.1"}
	bodySize, _ := msg.FullSize()
	head := messages.Head{5, 0, 0, 0, 0, 0, 0, bodySize}
	head.WriteTo(client.conn)
	msg.WriteTo(client.conn)
	timeChan := make(chan int32)
	go func() {
		for {
			client.read()
		}
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			client.sendTime(timeChan)
		}
	}()

	for {
		client.latencySec, client.latencyUsec = <-timeChan, <-timeChan
	}
}

func getDecoder(codecHeader messages.Codec) (func([]byte, messages.ServerSettings) ([]int16, error), error) {
	switch codecHeader.Codec {
	case messages.OPUS:
		var frameMs, middling float32 = 60, 1
		opusHeader := codecHeader.Payload.(*messages.OpusHeader)
		channels := opusHeader.Channels
		sampleRate := opusHeader.SampleRate
		middling = 1.0 / float32(channels)
		frameSize := float32(channels) * frameMs * float32(sampleRate) / 1000
		d, _ := opus.NewDecoder(int(sampleRate), int(channels))
		return func(data []byte, setting messages.ServerSettings) ([]int16, error) {
			pcm1 := make([]int16, int(frameSize))
			n, _ := d.Decode(data, pcm1)
			output := make([]int16, n)
			pcm1 = pcm1[:n*int(channels)]
			for i := 0; i < n; i++ {
				for y := 0; y < int(channels); y++ {
					ch := float32(pcm1[(i*int(channels))+y]) * middling * float32(setting.Volume) / 100
					output[i] += int16(ch)
				}
			}
			return output, nil
		}, nil
	}
	err1 := errors.New("Foobar")
	return nil, err1
}

func (client *SnapClient) read() {
	var buffer []byte
	head := messages.Head{}
	_, err := head.ReadFrom(client.conn)
	if err != nil {
		client.log.Error(err)
	}

	switch head.MsgType {
	case messages.CodecMsg:
		client.codecHeader = messages.Codec{}
		client.codecHeader.ReadFrom(client.conn)
		client.d, _ = getDecoder(client.codecHeader)
	case messages.WireChunkMsg:
		payload := messages.WireChunk{}
		_, err := payload.ReadFrom(client.conn)
		if err != nil {
			client.log.Error(err)
		}
		if client.d != nil {
			result, _ := client.d([]byte(payload.Payload), client.serverSetting)
			client.output <- result
		}
	case messages.ServerSettingMsg:
		payload := messages.ServerSettings{}
		_, err := payload.ReadFrom(client.conn)
		if err != nil {
			client.log.Error(err)
		} else {
			client.serverSetting = payload
		}
	case messages.TimeMsg:
		client.timeHead = head
		payload := messages.Time{}
		payload.ReadFrom(client.conn)
		client.timeMessage = payload
	case messages.HelloMsg:
	case messages.StreamTagMsg:
		var bodySize uint32
		buffer = make([]byte, 4)
		client.conn.Read(buffer)
		r := bytes.NewReader(buffer)
		binary.Read(r, binary.LittleEndian, &bodySize)
		payLoadBytes := make([]byte, bodySize)
		client.conn.Read(payLoadBytes)
		r = bytes.NewReader(payLoadBytes)
		binary.Read(r, binary.LittleEndian, &payLoadBytes)
		//payLoad := string(payLoadBytes)
		//client.stdLog.Println(payLoad)
	}
}

func (client *SnapClient) sendTime(c chan int32) {
	c <- client.timeHead.Received_sec - client.timeHead.Sent_sec
	c <- client.timeHead.Received_usec - client.timeHead.Sent_usec
	head := messages.Head{4, 0, client.timeHead.Id, 0, 0, 0, 0, 8}
	head.WriteTo(client.conn)
	time := messages.CreateTime()

	time.LatencySec = client.latencySec
	time.LatencyUsec = client.latencyUsec
	time.WriteTo(client.conn)
}
