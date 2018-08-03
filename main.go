package main

import (
	"errors"
	"github.com/spf13/viper"
	"gitlab.com/g00g/vkcli/tools"
	"log"
	"path/filepath"
)

func main() {
	Start()
}

func Start() {
	dataFolder := "./data"
	configFileName := "config.yaml"

	initConfig(filepath.Join(dataFolder, configFileName))
	//TODO use config.NewClient() to make requests to VK server
}

func initConfig(pathToConfig string) {
	err := tools.ReadConfig(pathToConfig)
	if err != nil {
		log.Fatal(err)
	}

	checkConfig()
}

func checkConfig() {
	if got := viper.Get("clientId"); got == nil || got == "" {
		log.Fatal(errors.New("config file does not contain clientId"))
	}
}
