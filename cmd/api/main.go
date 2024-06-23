package main

import (
	"backend/config"
	"backend/server"
	"flag"
	"log"

	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	var isProduction bool
    flag.BoolVar(&isProduction, "prod", false, "Çalışma ortamı")
	flag.Parse()
	app := server.NewApp(isProduction)

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
