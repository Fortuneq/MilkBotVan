package config

import (
	"github.com/go-playground/validator/v10"
	zerolog "github.com/rs/zerolog"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	Telegram           Telegram           `validate:"required"`
	RegistrationSystem RegistrationSystem `validate:"required"`
	Postgres           Postgres           `validate:"required"`
	RabbitMQ           RabbitMQ           `validate:"required"`
	Logger             Logger             `validate:"required"`
}

type Telegram struct {
	Error        map[string]BotConfig `validate:"required"`
	TronWallet   BotConfig            `validate:"required"`
	Treasury     BotConfig            `validate:"required"`
	Logging      BotConfig            `validate:"required"`
	VerifySystem VerifySystem         `validate:"required"`
}

type BotConfig struct {
	ChatID string `valdiate:"required"`
	Token  string `validate:"required"`
}

type VerifySystem struct {
	Host           string `validate:"required"`
	Port           string `validate:"required"`
	Token          string `validate:"required"`
	SupportChatID1 string `validate:"required"`
	SupportChatID2 string `validate:"required"`
	SupportChatID3 string `validate:"required"`
	SupportChatID4 string `validate:"required"`
}
type RegistrationSystem struct {
	UserTokenLifetime        time.Duration `validate:"required"`
	ClientTokenLifetime      time.Duration `validate:"required"`
	SellerAgentTokenLifetime time.Duration `validate:"required"`
}

type Postgres struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
	PGDriver string `validate:"required"`
	Settings struct {
		MaxOpenConns    int           `validate:"required,min=1"`
		ConnMaxLifetime time.Duration `validate:"required,min=1"`
		MaxIdleConns    int           `validate:"required,min=1"`
		ConnMaxIdleTime time.Duration `validate:"required,min=1"`
	}
}

type RabbitMQ struct {
	Host                     string `validate:"required"`
	Port                     string `validate:"required"`
	User                     string `validate:"required"`
	Password                 string `validate:"required"`
	ErrorExchangeName        string `validate:"required"`
	TronWalletExchangeName   string `validate:"required"`
	NotificationExchangeName string `validate:"required"`
	TreasuryExchangeName     string `validate:"required"`
	LoggingExchangeName      string `validate:"required"`
}

type Logger struct {
	Level *zerolog.Level `validate:"required"`
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	path := os.Getenv("CONFIG")

	if len(path) != 0 {
		v.AddConfigPath(path)
	}

	v.AddConfigPath("config")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.BindEnv("Postgres.Password", "POSTGRES_PASSWORD")
	v.BindEnv("RabbitMQ.Password", "RABBITMQ_PASSWORD")
	v.BindEnv("Telegram.VerifySystem.Token", "TELEGRAM_VERIFY_SYSTEM_TOKEN")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return nil, err
	}
	err = validator.New().Struct(c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
