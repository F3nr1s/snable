package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mumble struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Retries  int    `yaml:"retries"`
	} `yaml:"mumble"`
	Volume struct {
		Default	 uint	`yaml:"default"`
		Ducking	 uint	`yaml:"ducking"`
	} `yaml:"volume"`
	Snapcast struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		WebPort	string `yaml:"webport"`
		Retries int    `yaml:"retries"`
	} `yaml:"snapcast"`
}

func Loadconfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	if config.Mumble.Retries < 1 {
		config.Mumble.Retries = 3
	}

	if config.Snapcast.Retries < 1 {
		config.Snapcast.Retries = 3
	}

	if config.Volume.Default > 100 {
		config.Volume.Default = 100
	}

	if config.Volume.Ducking > 100 {
		config.Volume.Ducking = 100
	}

	return config, nil
}

func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
