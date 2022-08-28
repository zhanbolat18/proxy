package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhanbolat18/proxy/config"
	controller "github.com/zhanbolat18/proxy/deliveries/http"
	"github.com/zhanbolat18/proxy/internal/proxy/repositories/memory"
	httpUseCase "github.com/zhanbolat18/proxy/internal/proxy/usecases/http"
	"github.com/zhanbolat18/proxy/pkg/doer"
	"github.com/zhanbolat18/proxy/pkg/generator/uuid"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg := config.NewConfig()
	gengine := provideDependencies(cfg)
	srv := &http.Server{
		Addr:        cfg.Srv.Port,
		Handler:     gengine,
		ReadTimeout: cfg.Srv.ReadTimeout,
		ErrorLog:    log.Default(),
	}
	go startServer(srv)
	gracefulShutdown(srv, cfg.Srv)
}

func provideDependencies(cfg *config.Config) *gin.Engine {
	idGen := uuid.NewUUIDGenerator()
	httpDoer := httpClient(cfg.Client)
	repo := memory.NewMemoryRepository()
	uc := httpUseCase.NewHttpUsecase(httpDoer, idGen, repo)
	cntr := controller.NewProxyController(uc)

	gengine := gin.Default()
	gengine.GET("/", cntr.Proxy)
	return gengine
}

func httpClient(client *config.HttpClient) doer.Doer {
	return &http.Client{
		Timeout: client.Timeout,
	}
}

func startServer(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}
}

func gracefulShutdown(srv *http.Server, cfg *config.HttpServer) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	ctx, cf := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cf()
	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
		cf()
	}()
	<-ctx.Done()
}
