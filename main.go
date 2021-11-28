package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/f3nr1s/snable/mumble"
	"github.com/f3nr1s/snable/snapcast"
)

func main() {
	configPath, err := ParseFlags()
	if err != nil {
		panic(err)
	}
	config, err := Loadconfig(configPath)
	if err != nil {
		panic(err)
	}

	fob, _ := net.Interfaces()
	mac := "11:22:33:44:55:66"
	for _, ifa := range fob {
		mac = ifa.HardwareAddr.String()
		if mac != "" {
			break
		}
	}
	ch := make(chan []int16)
	i := 0
	snapclient, err := snapcast.Create(config.Snapcast.Host, config.Snapcast.Port, mac, ch)
	for err != nil {
		i++
		fmt.Println(err)
		if i >= config.Snapcast.Retries {
			fmt.Printf("Reach retry amount %d, stopping\n", config.Snapcast.Retries)
			os.Exit(1)
		}
		fmt.Println("Trying again in one second")
		time.Sleep(1 * time.Second)

		snapclient, err = snapcast.Create(config.Snapcast.Host, config.Snapcast.Port, mac, ch)
	}
	defer snapclient.Close()
	go snapclient.Initialize()

	controller, _ := snapcast.CreateSnapControl(config.Snapcast.Host, config.Snapcast.WebPort, mac, config.Volume.Default, config.Volume.Ducking)

	mumbleClient, err := mumble.Create(
		config.Mumble.Host,
		config.Mumble.Port,
		config.Mumble.Username,
		config.Mumble.Password,
		controller)
	i = 0
	for err != nil {
		i++
		fmt.Println(err)
		if i >= config.Mumble.Retries {
			fmt.Printf("Reach retry amount %d, stopping\n", config.Mumble.Retries)
			os.Exit(1)
		}
		fmt.Println("Trying again in one second")
		time.Sleep(1 * time.Second)

		mumbleClient, err = mumble.Create(
			config.Mumble.Host,
			config.Mumble.Port,
			config.Mumble.Username,
			config.Mumble.Password,
			controller)
	}
	defer mumbleClient.Conn.Close()

	mumble.Start(mumbleClient, ch)
}
