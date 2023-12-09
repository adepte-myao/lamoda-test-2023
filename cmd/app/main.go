package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adepte-myao/lamoda-test-2023/configs"
	loggers "github.com/adepte-myao/lamoda-test-2023/internal/pkg/logger"
	"github.com/adepte-myao/lamoda-test-2023/internal/pkg/server"
)

func main() {
	cfg, err := configs.LoadDefault()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := loggers.NewZap(cfg)
	if err != nil {
		log.Fatal(err)
	}

	engine := gin.New()

	engine.Use(gin.Recovery())

	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, map[string]string{"info": "pong"})
	})

	httpServer := server.New(cfg, engine.Handler(), logger)

	err = httpServer.Run()
	if err != nil {
		logger.Error(err)
	}
}
