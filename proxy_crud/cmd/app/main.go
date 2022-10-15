package main

import (
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	_ "proxy_crud/docs"
	"proxy_crud/internal/app"
	"proxy_crud/internal/config"
	"proxy_crud/pkg/logging"
)

// @title           Proxy Crud Service
// @version         1.0
// @description     CRUD.

// @host      localhost:10000
// @BasePath  /api/proxy_crud/v1

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
