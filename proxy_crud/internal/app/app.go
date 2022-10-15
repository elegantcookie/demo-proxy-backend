package app

import (
	"context"
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
	"proxy_crud/pkg/client/postgresql"
	"proxy_crud/pkg/logging"
	"proxy_crud/pkg/metrics"
	"time"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *gin.Engine
	httpServer *http.Server
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

	storage := db.NewStorage(psqlClient, logger)
	fmt.Printf("%v", storage)

	service, _ := proxy.NewService(storage, logger)

	//p := []proxy.Proxy{
	//	proxy.NewProxy(proxy.CreateProxyDTO{
	//		Ip:         "1.2.3.4",
	//		Port:       10,
	//		ExternalIP: "1.2.3.4",
	//		Country:    "Russia",
	//	}),
	//	proxy.NewProxy(proxy.CreateProxyDTO{
	//		Ip:         "1.2.3.4",
	//		Port:       10,
	//		ExternalIP: "1.2.3.4",
	//		Country:    "Russia",
	//	}),
	//	proxy.NewProxy(proxy.CreateProxyDTO{
	//		Ip:         "1.2.3.4",
	//		Port:       10,
	//		ExternalIP: "1.2.3.4",
	//		Country:    "Russia",
	//	}),
	//	proxy.NewProxy(proxy.CreateProxyDTO{
	//		Ip:         "1.2.3.4",
	//		Port:       10,
	//		ExternalIP: "1.2.3.4",
	//		Country:    "Russia",
	//	}),
	//	proxy.NewProxy(proxy.CreateProxyDTO{
	//		Ip:         "1.2.3.4",
	//		Port:       10,
	//		ExternalIP: "1.2.3.4",
	//		Country:    "Russia",
	//	}),
	//}
	//for i := 0; i < 5; i++ {
	//	id := uuid.NewV4()
	//	p[i].ID = id.String()
	//}
	//
	//err = storage.Insert(context.TODO(), p)
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//p := proxy.Proxy{
	//	Ip:         "1.2.3.4",
	//	Port:       10,
	//	ExternalIP: "1.2.3.4",
	//	CreatedAt:  time.Now(),
	//	Country:    "Zimbabve",
	//}
	//
	//id, err := storage.InsertOne(context.TODO(), p)
	//if err != nil {
	//	log.Fatal("%v", err)
	//}
	//fmt.Println(id)
	//p, err := storage.FindById(context.TODO(), "88edcf8d-110b-447b-8252-3becbbd86626")
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//fmt.Println(p)

	//proxies, err := storage.FindAll(context.TODO())
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//fmt.Println(proxies)

	//redisClient, err := redis.NewClient()
	//if err != nil {
	//	panic(err)
	//}
	//logger.Infof("redis client: %+v\nis null: %v", redisClient, redisClient == nil)
	//cache := rcache.NewCache(redisClient, logger)
	//storage := db.NewStorage(mongodbClient, "user_service", logger)
	//service, err := user.NewService(storage, cache, *logger)
	//if err != nil {
	//	panic(err)
	//}

	proxyHandler := proxy.Handler{
		Logger:       *logger,
		ProxyService: service,
	}
	proxyHandler.Register(v1)

	return App{
		cfg,
		logger,
		router,
		nil,
	}, nil
}

func (a *App) Run() {
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
