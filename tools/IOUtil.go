package tools

import (
	"os"
	"github.com/spf13/viper"
	"errors"
	"fmt"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func ReadConfig(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read config file: ", err))
	}
	return nil
}
