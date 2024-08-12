package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Default Config

func Load(file ...string) error {
	if len(file) > 0 && len(file[0]) > 0 {
		viper.SetConfigFile(file[0])
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&Default); err != nil {
		return fmt.Errorf("parser config err: %w", err)
	}

	return nil
}
