package main

import (
	"fmt"

	"github.com/spf13/viper"
)

//InitConfig - This method will initialize a Viper configuration object with our config.yml file
func InitConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file was not found: %s", err.Error()))
		} else {
			panic(fmt.Errorf("Config file was found but other error is present: %s", err.Error()))
		}
	}
	return v, err
}
