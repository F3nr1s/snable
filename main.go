package main

import (
	"net"
	"os"
	"time"

	"github.com/f3nr1s/snable/mumble"
	"github.com/f3nr1s/snable/snapcast"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath, err := ParseFlags()
	if err != nil {
		logrus.Fatal("Error")
	}
	config, err := Loadconfig(configPath)
	if err != nil {
		logrus.Fatal(err)
	}
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(logrus.Level(config.Debug.Level))
	if config.Debug.LogFile != "" {
		file, err := os.OpenFile(config.Debug.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.WithField("filename", config.Debug.LogFile).Fatal("Can't open Log")
		}
		log.SetOutput(file)
	}

	fob, _ := net.Interfaces()
	mac := "11:22:33:44:55:66"
	for _, ifa := range fob {
		mac = ifa.HardwareAddr.String()
		if mac != "" {
			break
		}
	}
	log.WithField("mac", mac).Debug("Found mac address")
	ch := make(chan []int16)
	i := 0
	snapclient, err := snapcast.Create(
		config.Snapcast.Host,
		config.Snapcast.Port,
		mac,
		ch,
		log)
	for err != nil {
		i++
		log.Error(err)
		if i >= config.Snapcast.Retries {
			log.WithField("amount", config.Snapcast.Retries).Fatal("Reach retry amount, stopping")
		}
		log.Debug("Trying again in one second")
		time.Sleep(1 * time.Second)

		snapclient, err = snapcast.Create(
			config.Snapcast.Host,
			config.Snapcast.Port,
			mac,
			ch,
			log)
	}
	defer snapclient.Close()
	go snapclient.Initialize()

	controller, _ := snapcast.CreateSnapControl(
		config.Snapcast.Host,
		config.Snapcast.WebPort,
		mac,
		config.Volume.Default,
		config.Volume.Ducking)

	mumbleClient, err := mumble.Create(
		config.Mumble.Host,
		config.Mumble.Port,
		config.Mumble.Username,
		config.Mumble.Password,
		controller)
	i = 0
	for err != nil {
		i++
		log.Warn(err)
		if i >= config.Mumble.Retries {
			log.WithField("amount", config.Mumble.Retries).Fatal("Reach retry amount, stopping")
		}
		log.Debug("Trying again in one second")
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
