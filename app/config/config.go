package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server        Server
	Logger        Logger
	OpenTelemetry OpenTelemetry
	Redis         Redis
	TTLUUID       int `validate:"required"`
}

var C *Config

func GetConfig() *Config {
	return C
}

type Server struct {
	AppVersion                  string `validate:"required"`
	Host                        string `validate:"required"`
	GRPCPort                    string `validate:"required"`
	HTTPPort                    string `validate:"required"`
	ShowUnknownErrorsInResponse bool
	IPHeader                    string `validate:"required"`
}

type Logger struct {
	Level          string `validate:"required"`
	SkipFrameCount int
	InFile         bool
	FilePath       string
	InTG           bool
	TGLevel        string `validate:"required"`
	ChatID         int64
	TGToken        string
	AlertUsers     []string
}

type OpenTelemetry struct {
	Host        string `validate:"required"`
	ServiceName string `validate:"required"`
}

type Redis struct {
	Host               string `validate:"required"`
	Port               string `validate:"required"`
	MinIdleConns       int    `validate:"required"`
	PoolSize           int    `validate:"required"`
	PoolTimeout        int    `validate:"required"`
	Password           string `validate:"required"`
	UseCertificates    bool
	InsecureSkipVerify bool
	CertificatesPaths  struct {
		Cert string
		Key  string
		Ca   string
	}
	DB int
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(fmt.Sprintf("./%s", ConfigPath))
	v.SetConfigName(ConfigFileName)
	v.SetConfigType(ConfigExtension)
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
