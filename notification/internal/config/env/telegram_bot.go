package env

import "github.com/caarlos0/env/v11"

type telegramBotEnvConfig struct {
	Token string `env:"TELEGRAM_BOT_TOKEN,required"`
}

type telegramBotConfig struct {
	raw telegramBotEnvConfig
}

func NewTelegramBotConfig() (*telegramBotConfig, error) {
	var raw telegramBotEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &telegramBotConfig{raw: raw}, nil
}

func (tbc *telegramBotConfig) Token() string {
	return tbc.raw.Token
}
