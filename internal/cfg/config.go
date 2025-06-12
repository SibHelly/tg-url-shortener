package cfg

import (
	"log"

	"github.com/spf13/viper"
)

type Cfg struct {
	Token string
}

func LoadConfig() *Cfg {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	return &Cfg{
		Token: viper.GetString("BOT_TOKEN"),
	}
}
