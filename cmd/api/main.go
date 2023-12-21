package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"ail-test/cmd/api/config"
	"ail-test/cmd/api/server"
)

func main() {
	cfg := config.Config{}
	log := logrus.StandardLogger()
	envconfig.MustProcess("API", &cfg)
	log.Info("starting server : ", cfg.Name)
	app := server.Handler(&cfg)
	startServerString := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	err := app.Listen(startServerString)
	if err != nil {
		panic(err)
	}
}
