package main

import (
	"log"
	"proxy_crud/internal/app"
	"proxy_crud/internal/config"
	"proxy_crud/pkg/logging"
)

func main() {
	log.Print("config initialization")
	cfg := config.GetConfig()

	log.Printf("logging initialized.")

	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("running Application")
	a.Run()
}
