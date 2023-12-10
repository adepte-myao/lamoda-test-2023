package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/adepte-myao/lamoda-test-2023/configs"
	"github.com/adepte-myao/lamoda-test-2023/docs"
	loggers "github.com/adepte-myao/lamoda-test-2023/internal/pkg/logger"
	"github.com/adepte-myao/lamoda-test-2023/internal/pkg/postgres"
	"github.com/adepte-myao/lamoda-test-2023/internal/pkg/server"
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/services"
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/handlers"
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/repositories"
)

// @title Reservation microservice
// @version 1.0

func main() {
	cfg, err := configs.LoadDefault()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := loggers.NewZap(cfg)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New()

	dbCfg := cfg.Database
	postgresDB, cancelDB, err := postgres.NewDB(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DB))
	if err != nil {
		logger.Fatal(err)
	}

	defer cancelDB(postgresDB)

	storehouseRepo := repositories.NewPostgresStorehouse(postgresDB)
	itemRepo := repositories.NewPostgresItem(postgresDB)
	reservationRepo := repositories.NewPostgresReservation(postgresDB)

	service := services.New(storehouseRepo, itemRepo, reservationRepo)

	handler := handlers.NewReservationHandler(service, validate)

	engine := gin.New()
	engine.Use(gin.Recovery())

	if cfg.Logger.Level == "debug" {
		engine.Use(loggers.SendErrorsToClient)
	}

	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, map[string]string{"info": "pong"})
	})

	docs.SwaggerInfo.BasePath = "/"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.POST("/reserve", handler.Reserve)
	engine.POST("/release", handler.Release)
	engine.GET("/get-unreserved-items", handler.GetUnreserved)

	httpServer := server.New(cfg, engine.Handler(), logger)

	err = httpServer.Run()
	if err != nil {
		logger.Error(err)
	}
}
