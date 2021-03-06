package tools

import (
	"testing"
	//"gitlab.com/g00g/vk-cli/tools"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

var dataFolder = "./data"

func TestReadConfig(t *testing.T) {
	conf := "config.yaml"
	if !FileExists(dataFolder) {
		os.MkdirAll(dataFolder, os.ModePerm)
	}
	defer os.RemoveAll(dataFolder)

	err := ioutil.WriteFile(filepath.Join(dataFolder, conf), []byte("client_id: \"0123456789\""),
		0600)
	if err != nil {
		t.Fatalf("Cannot write a config file: %v", err)
	}

	err = ReadConfig(filepath.Join(dataFolder, conf))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := viper.Get("client_id"), "0123456789"; got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}
