package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv              string `mapstructure:"APP_ENV"`
	AppPort             string `mapstructure:"APP_PORT"`
	VoucherifyId        string `mapstructure:"VOUCHERIFY_ID"`
	VoucherifySecretKey string `mapstructure:"VOUCHERIFY_SECRET_KEY"`
}

var log zerolog.Logger

func InitEnv() (*Env, zerolog.Logger) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log = zerolog.New(output).
		Level(zerolog.Level(zerolog.InfoLevel)).
		With().Timestamp().
		Logger()

	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Msg("Error reading .env file")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal().Msg("Error unmarshalling .env file")
	}

	if env.AppEnv == "DEV" {
		log.Info().Msgf("DEV mode:%s", env.AppPort)
	}

	return &env, log
}
