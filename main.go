package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab.com/g00g/vkcli/auth"
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

	initFromConfigFile(filepath.Join(dataFolder, configFileName))
	//TODO use config.NewClient() to make requests to VK server
}

func initFromConfigFile(pathToConfig string) (clientId string, tokenUrl string) {
	err := tools.ReadConfig(pathToConfig)
	if err != nil {
		log.Fatal(err)
	}

	clientId = viper.GetString("clientId")
	if clientId == "" {
		log.Fatal(fmt.Errorf("confing file %v has error: clientId not found", pathToConfig))
	}

	tokenUrl = viper.GetString("tokenUrl")
	if tokenUrl == "" {
		requestToken(clientId, pathToConfig)
	}
	return
}

func requestToken(clientId string, pathToConfig string) {
	config := auth.NewConfig(clientId)
	log.Fatal(fmt.Errorf("tokenUrl: not found in confing file. Looks like this is the first time you run VK-Cli. "+
		"Go to the following URL, log in, allow access, than copy url from the adress line. Add copied URL into %v config file"+
		" as a new line.\nFor example tokenUrl: \"url goes here\"\n%v", pathToConfig, config.AuthCodeURL("", config.DefaultOptions())))
}
