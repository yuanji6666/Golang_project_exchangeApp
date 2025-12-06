package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn string
		MaxIdlesConns int
		MaxOpenConns int
	}
}

var Appconfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file : %v", err)
	}

	Appconfig = &Config{}

	err := viper.Unmarshal(Appconfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v",err)
	}

	initDB()

}

