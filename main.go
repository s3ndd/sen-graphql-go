package main

import (
	"net/http"
	"time"

	internalConfig "github.com/s3ndd/sen-graphql-go/internal/config"
	"github.com/s3ndd/sen-graphql-go/router"

	"github.com/gin-gonic/gin"
	"github.com/s3ndd/sen-go/config"
	"github.com/s3ndd/sen-go/log"
)

func main() {
	config.Load()
	logger := log.NewLogger(internalConfig.Logger())
	logger.Info("Starting...")
	if env := config.String("ENV", "dev"); env == "prod" || env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	routerInit := router.InitRouter()

	errs := make(chan error, 1)
	go func() {
		server := &http.Server{
			Addr:         ":" + config.String("PORT", "60000"),
			Handler:      routerInit,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 60,
		}
		errs <- server.ListenAndServe()
	}()
	logger.WithError(<-errs).Info("terminated")
}
