package app

import (
	"bytes"
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
	"proxy_crud/internal/config"
	"proxy_crud/internal/proxy"
	"proxy_crud/internal/proxy/db"
	"proxy_crud/internal/proxy/service"
	proxy_group "proxy_crud/internal/proxy_group"
	pgdb "proxy_crud/internal/proxy_group/db"
	pgservice "proxy_crud/internal/proxy_group/service"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/client/kafka"
	"proxy_crud/pkg/client/postgresql"
	"proxy_crud/pkg/logging"
	"proxy_crud/pkg/metrics"
	"time"
)

type App struct {
	cfg           *config.Config
	logger        *logging.Logger
	kafkaProducer kafka.IProducer
	service       service.Service
	router        *gin.Engine
	httpServer    *http.Server
}

const (
	statusProcessing = 1
)

func startProducer(ctx context.Context, logger *logging.Logger, producer kafka.IProducer, service service.Service) {

	defer producer.Close()
	var options filter.Options
	sopts := filter.SOptions{
		Field: "checked_at",
		Order: "ASC",
	}
	initFields := make([]filter.Field, 0)
	limit := 10000
	fopts := filter.NewFOptions(true, limit, 1, initFields)
	options.SortOptions = sopts
	options.FilterOptions = fopts
	options.ValidateIntAndAdd("processing_status", "0", filter.OperatorEqual)
	for true {
		proxies, err := service.GetAll(ctx, options)
		if len(proxies) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			return
		}
		for i := 0; i < len(proxies); i++ {
			i := i
			go func() {
				proxies[i].ProcessingStatus = 1
				err := service.UpdateProxyStatus(ctx, proxies[i].ID, statusProcessing)
				if err != nil {
					logger.Errorln(err)
				}
				var network bytes.Buffer        // Stand-in for a network connection
				enc := gob.NewEncoder(&network) // Will write to network.

				err = enc.Encode(proxies[i])
				if err != nil {
					logger.Errorln(err)
				}

				err = producer.Produce(ctx, []byte(proxies[i].ID), network.Bytes())
				if err != nil {
					logger.Errorln(err)
				}
			}()
		}
		time.Sleep(1 * time.Second)
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

	psqlClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	proxyStorage := db.NewStorage(psqlClient, logger)
	fmt.Printf("%v", proxyStorage)

	proxyGroupStorage := pgdb.NewStorage(psqlClient, logger)

	proxyGroupService, _ := pgservice.NewService(proxyGroupStorage, logger)
	proxyService, _ := service.NewService(proxyStorage, proxyGroupStorage, logger)

	proxyHandler := proxy.Handler{
		Logger:       *logger,
		ProxyService: proxyService,
	}
	proxyGroupHandler := proxy_group.Handler{
		Logger:            *logger,
		ProxyGroupService: proxyGroupService,
	}
	proxyPath := v1.Group("/proxy")
	proxyHandler.Register(proxyPath)

	proxyGroupPath := v1.Group("/proxy_group")
	proxyGroupHandler.Register(proxyGroupPath)

	kafkaProducer := kafka.NewClient(context.TODO(), cfg.Kafka.URL, cfg.Kafka.Topic)

	return App{
		cfg,
		logger,
		kafkaProducer,
		proxyService,
		router,
		nil,
	}, nil
}

func (a *App) Run() {
	go startProducer(context.Background(), a.logger, a.kafkaProducer, a.service)
	a.startHTTP()
}

func (a *App) startHTTP() {
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
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}

}
