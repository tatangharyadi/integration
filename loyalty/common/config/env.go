package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv              string `mapstructure:"APP_ENV"`
	VoucherifyId        string `mapstructure:"VOUCHERIFY_ID"`
	VoucherifySecretKey string `mapstructure:"VOUCHERIFY_SECRET_KEY"`
}

func InitEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading .env file")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Error unmarshalling .env file")
	}

	if env.AppEnv == "Dev" {
		log.Println("Running in Dev mode")
	}

	return &env
}
