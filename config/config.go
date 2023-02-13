package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig config

type config struct {
	Currencies   []string `mapstructure:"CURRENCIES"`
	Interval     int      `mapstructure:"INTERVAL"`
	BaseCurrency string   `mapstructure:"BASECURRENCY"`
}

func InitConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
}
