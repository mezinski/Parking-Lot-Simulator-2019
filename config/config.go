package config

import (
	"fmt"

	"github.com/spf13/viper"
)

//InitConfig ...
func InitConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./config")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file was not found: %s", err.Error()))
		} else {
			panic(fmt.Errorf("Config file was found but other error is present: %s", err.Error()))
		}
	}
	return v, err
}
