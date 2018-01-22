package main

import (
	"bytes"
	"fmt"
	"time"

	"os"

	"github.com/BurntSushi/toml"
	"github.com/vidmed/logger"
)

var configInstance *TomlConfig

// TomlConfig represents a config
type TomlConfig struct {
	Main Main
}

// Main represent a main section of the TomlConfig
type Main struct {
	LogLevel        uint8
	ListenStr       string
	FlushPeriod     uint
	MaxTransactions uint
	FlushFile       string
}

// GetConfig returns application config
func GetConfig() *TomlConfig {
	return configInstance
}

// NewConfig creates new application config with given .toml file
func NewConfig(file string) (*TomlConfig, error) {
	configInstance = &TomlConfig{}

	if _, err := toml.DecodeFile(file, configInstance); err != nil {
		return nil, err
	}
	dump(configInstance)

	// check required fields
	// Main
	if configInstance.Main.ListenStr == "" {
		logger.Get().Fatalln("Main.ListenStr must be specified. Check your Config file")
	}
	if configInstance.Main.FlushPeriod == 0 {
		logger.Get().Fatalln("Main.FlushPeriod must be specified. Check your Config file")
	}
	if configInstance.Main.MaxTransactions == 0 {
		logger.Get().Fatalln("Main.MaxTransactions must be specified. Check your Config file")
	}
	if configInstance.Main.FlushFile == "" {
		logger.Get().Fatalln("Main.FlushFile must be specified. Check your Config file")
	}

	f, err := os.OpenFile(configInstance.Main.FlushFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Get().Fatalf("Main.FlushFile check error: %s", err.Error())
	}
	f.Close()

	return configInstance, nil
}

func dump(cfg *TomlConfig) {
	var buffer bytes.Buffer
	e := toml.NewEncoder(&buffer)
	err := e.Encode(cfg)
	if err != nil {
		logger.Get().Fatal(err)
	}

	fmt.Println(
		time.Now().UTC(),
		"\n---------------------Sevice started with config:\n",
		buffer.String(),
		"\n---------------------")
}
