package mumble

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/f3nr1s/snable/snapcast"

	"github.com/bep/debounce"
	cmap "github.com/orcaman/concurrent-map"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleutil"
	_ "layeh.com/gumble/opus"
)

type audioHandler struct {
	counter       cmap.ConcurrentMap
	volController *snapcast.SnapControl
}

func (a *audioHandler) OnAudioStream(e *gumble.AudioStreamEvent) {
	if a.counter == nil {
		a.counter = cmap.New()
	}
	f := func() {
		a.counter.Remove(e.User.Name)
		if !(a.counter.Count() > 0) {
			a.volController.SetDuckingState(false)
		}
	}
	debounced := debounce.New(500 * time.Millisecond)
	go func() {
		for range e.C {
			a.counter.SetIfAbsent(e.User.Name, true)
			a.volController.SetDuckingState(true)
			debounced(f)
		}
	}()
}

func sendHelp(m *gumble.TextMessageEvent) {
	helpMess := []string{
		"<p><b>Help:</b><br/>-----</p>",
		"<p><b>Volume Control</b></p>",
		"<p><b>.vol [number]:</b><br/>Gets/Sets the volume, example: .vol 50<br/>",
		"<b>.duck [number]:</b><br/>Gets/Sets the ducking vol, example. .duck 20</p>",
		"<p>-----</p>"}
	message := gumble.TextMessage{}
	message.Users = append(message.Users, m.Sender)
	message.Message = strings.Join(helpMess, "\n")
	m.Client.Send(&message)
}

func sendUnknownCall(m *gumble.TextMessageEvent, c string) {
	message := gumble.TextMessage{}
	message.Users = append(message.Users, m.Sender)
	message.Message = "Unknown command: " + c
	m.Client.Send(&message)
}

func handleVol(m *gumble.TextMessageEvent, s []string, volController *snapcast.SnapControl) {
	message := gumble.TextMessage{}
	message.Users = append(message.Users, m.Sender)
	if len(s) > 0 {
		if n, err := strconv.ParseUint(s[0], 10, 8); err != nil && n <= 100 {
			fmt.Println(err)
		} else {
			volController.SetVolume(uint(n))
		}
	}
	message.Message = fmt.Sprintf("Current volume: %d", volController.GetVolume())
	m.Client.Send(&message)
}

func handleDuck(m *gumble.TextMessageEvent, s []string, volController *snapcast.SnapControl) {
	message := gumble.TextMessage{}
	message.Users = append(message.Users, m.Sender)
	if len(s) > 0 {
		if n, err := strconv.ParseUint(s[0], 10, 8); err != nil && n <= 100 {
			fmt.Println(err)
		} else {
			volController.SetDucking(uint(n))
		}
	}

	message.Message = fmt.Sprintf("Current duck volume: %d", volController.GetDucking())
	m.Client.Send(&message)
}

func Create(host, port, username, password string, volController *snapcast.SnapControl) (*gumble.Client, error) {
	config := gumble.NewConfig()
	config.Username = username
	config.Password = password
	config.AudioInterval, _ = time.ParseDuration("10ms")

	audioHandler := audioHandler{volController: volController}
	config.AttachAudio(&audioHandler)
	config.Attach(gumbleutil.Listener{
		TextMessage: func(e *gumble.TextMessageEvent) {
			if e.Message[0:1] == "." {
				s := strings.Split(e.Message, " ")
				switch s[0] {
				case ".help":
					sendHelp(e)
				case ".vol":
					handleVol(e, s[1:], volController)
				case ".duck":
					handleDuck(e, s[1:], volController)
				default:
					sendUnknownCall(e, s[0])
				}
			}
		},
	})
	client, err := gumble.Dial(host+":"+port, config)
	return client, err
}

func Start(client *gumble.Client, ch chan []int16) {
	config := client.Config
	framesize := config.AudioFrameSize()
	outgoing := client.AudioOutgoing()
	defer close(outgoing)
	var channelBuffer []int16
	outBuffer := make([]int16, 0, framesize)
	for {
		channelBuffer = <-ch
		for _, element := range channelBuffer {
			outBuffer = append(outBuffer, element)
			if len(outBuffer) == framesize {
				outgoing <- gumble.AudioBuffer(outBuffer)
				outBuffer = make([]int16, 0, framesize)
			}
		}
	}
}
