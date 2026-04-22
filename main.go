package main

import (
	"log"

	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	// TODO:
	// 1. pass config file path and redirect based on config.root

	config := loadConfig("config.toml")

	var conf Config

	_, err := toml.Decode(config, &conf)
	if err != nil {
		log.Fatal(err)
	}

	conf.createSession()
	conf.attachSession()
}

func loadConfig(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
