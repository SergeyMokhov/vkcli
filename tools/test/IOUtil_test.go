package test

import (
	"testing"
	"github.com/SergeyMokhov/vkcli/tools"
	"path/filepath"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

var dataFolder = "./data"

func TestReadConfig(t *testing.T) {
	conf := "config.yaml"
	if !tools.FileExists(dataFolder) {
		os.MkdirAll(dataFolder, os.ModePerm)
	}
	defer os.RemoveAll(dataFolder)

	err := ioutil.WriteFile(filepath.Join(dataFolder, conf), []byte("client_id: \"0123456789\""),
		0600)
	if err != nil {
		t.Fatalf("Cannot write a config file: %v", err)
	}

	err = tools.ReadConfig(filepath.Join(dataFolder, conf))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := viper.Get("client_id"), "0123456789"; got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}
