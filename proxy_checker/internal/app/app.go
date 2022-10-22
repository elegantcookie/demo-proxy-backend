package app

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"proxy_checker/internal/config"
	"proxy_checker/internal/proxy"
	"proxy_checker/pkg/client/kafka"
	"proxy_checker/pkg/logging"
	"proxy_checker/pkg/metrics"
	"strings"
	"time"
)

type App struct {
	cfg           *config.Config
	logger        *logging.Logger
	kafkaConsumer kafka.IConsumer
	proxyService  proxy.Service
	router        *gin.Engine
	httpServer    *http.Server
}

func startConsumer(ctx context.Context, logger *logging.Logger, consumer kafka.IConsumer, service proxy.Service) {
	logger.Info("consumer started...")
	for {
		message, err := consumer.ReadMessage(ctx)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		reader := strings.NewReader(string(message.Value))
		dec := gob.NewDecoder(reader) // Will write to network.

		var pr proxy.Proxy
		err = dec.Decode(&pr)
		if err != nil {
			logger.Errorln(err)
		}
		go func() {
			checkedProxy, err := service.Check(ctx, pr)
			if err != nil {
				logger.Errorln(err)
			}
			err = service.FetchChanges(ctx, checkedProxy)
			if err != nil {
				logger.Errorln(err)
			}
		}()

		// 200 messages a second
		time.Sleep(5 * time.Millisecond)
	}
}

func NewApp(cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Println("router initializing")

	// router w/o logger and recovery middleware
	router := gin.New()

	// add recovery middleware
	// recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	v1 := router.Group("/api/proxy_crud/v1")
	logger.Println("swagger docs initialization")
	// use wrappers to make it compatible with default http.HandleFunc
	v1.Handle(http.MethodGet, "/docs", gin.WrapH(http.RedirectHandler("/api/proxy_crud/v1/docs/index.html", http.StatusMovedPermanently)))
	v1.Handle(http.MethodGet, "/docs/*any", gin.WrapH(httpSwagger.WrapHandler))

	logger.Println("heartbeat metric initializing")
	metricHandler := metrics.Handler{}
	metricHandler.Register(v1)

	kafkaConsumer := kafka.NewClient(context.TODO(), cfg.Kafka.URL, cfg.Kafka.Topic, cfg.Kafka.GroupID)
	service, _ := proxy.NewService(logger)

	return App{
		cfg,
		logger,
		kafkaConsumer,
		service,
		router,
		nil,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			go startConsumer(ctx, a.logger, a.kafkaConsumer, a.proxyService)
			a.startHTTP(ctx)
		}
	}

}

func (a *App) startHTTP(ctx context.Context) {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == config.ListenTypeSock {
		appDir, err := filepath.Abs(os.Args[0])
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"https://localhost:3000", "https://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Authorization", "Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Access-Token", "Refresh-Token", "Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("application completely initialized and started")
	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		a.logger.Fatal(err)
	}

}
