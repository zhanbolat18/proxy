package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Srv    *HttpServer
	Client *HttpClient
}

type HttpServer struct {
	Port            string
	ReadTimeout     time.Duration
	ShutdownTimeout time.Duration
}
type HttpClient struct {
	Timeout time.Duration
}

func NewConfig() *Config {
	vpr := viper.New()
	vpr.AutomaticEnv()
	vpr.SetDefault("APP_PORT", ":8080")
	vpr.SetDefault("SRV_READ_TIMEOUT", time.Second*5)
	vpr.SetDefault("CLIENT_TIMEOUT", time.Second*15)
	vpr.SetDefault("SHUTDOWN_TIMEOUT", time.Second*30)

	return &Config{
		Srv: &HttpServer{
			Port:            vpr.GetString("APP_PORT"),
			ReadTimeout:     vpr.GetDuration("READ_TIMEOUT"),
			ShutdownTimeout: vpr.GetDuration("SHUTDOWN_TIMEOUT"),
		},
		Client: &HttpClient{
			Timeout: vpr.GetDuration("CLIENT_TIMEOUT"),
		},
	}
}
