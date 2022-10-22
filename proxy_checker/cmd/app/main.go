package main

import (
	"context"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	"os/signal"
	"proxy_checker/internal/app"
	"proxy_checker/internal/config"
	"proxy_checker/pkg/logging"
	"syscall"
)

// @title           Proxy Checker Service
// @version         1.0
// @description     Checks proxy from crud.

// @host      localhost:10000
// @BasePath  /api/checker/v1

func main() {
	log.Print("config initialization")
	cfg := config.GetConfig()

	log.Printf("logging initialized.")

	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("running application")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel
	err = a.Run(ctx)
	if err != nil {
		log.Println("stopping application")
	}
}
