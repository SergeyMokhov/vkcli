package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab.com/g00g/vk-cli/auth"
	"gitlab.com/g00g/vk-cli/client"
	"gitlab.com/g00g/vk-cli/tools"
	"log"
	"path/filepath"
	"strconv"
)

func main() {
	Start()
}

func Start() {
	dataFolder := "./data"
	configFileName := "config.yaml"

	tokenString := initFromConfigFile(filepath.Join(dataFolder, configFileName))

	token, err := auth.ParseUrlString(tokenString)
	if err != nil {
		log.Fatal(fmt.Errorf("error parsing token: %v", err))
	}

	vk := client.NewVk(token)
	//vk.ListFriends()
	vk.RemoveDeletedFriends()
	//vk.AddFriend(155633421)
	//vk.AddFriend(28421522)
	//vk.DeleteFriend(2916112)
	//vk.DeleteFriend(28421522)
}

func initFromConfigFile(pathToConfig string) (tokenUrl string) {
	err := tools.ReadConfig(pathToConfig)
	if err != nil {
		log.Fatal(err)
	}

	clientId := viper.GetString("clientId")
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
	config := auth.NewVkConfig(clientId)
	config.Scopes = []string{strconv.FormatInt(int64(auth.Friends()), 10)}
	log.Fatal(fmt.Errorf("tokenUrl: not found in confing file. Looks like this is the first time you run Vk-Cli. "+
		"Go to the following URL, log in, allow access, than copy url from the adress line. Add copied URL into %v config file"+
		" as a new line.\nFor example tokenUrl: \"url goes here\"\n%v", pathToConfig, config.AuthCodeURL("", config.DefaultOptions())))
}
