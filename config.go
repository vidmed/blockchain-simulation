package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/vidmed/logger"
	"time"
)

var config_instance *TomlConfig

type TomlConfig struct {
	Main MainConfig
}

type MainConfig struct {
	ListenStr   string
	LogLevel    int
	FlushPeriod int
	FlushFile   string
}

func GetConfig() *TomlConfig {
	return config_instance
}

func NewConfig(file string) (*TomlConfig, error) {
	config_instance = &TomlConfig{}

	if _, err := toml.DecodeFile(file, config_instance); err != nil {
		return nil, err
	}
	dump(config_instance)

	// check required fields
	// Main
	if config_instance.Main.ListenStr == "" {
		logger.Get().Fatalln("Main.ListenStr must be specified. Check your Config file")
	}
	if config_instance.Main.FlushPeriod == 0 {
		logger.Get().Fatalln("Main.FlushPeriod must be specified. Check your Config file")
	}
	if config_instance.Main.FlushFile == "" {
		logger.Get().Fatalln("Main.FlushFile must be specified. Check your Config file")
	}
	// todo check if file writable

	return config_instance, nil
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
		"\n---------------------\n\n\n\n")
}
