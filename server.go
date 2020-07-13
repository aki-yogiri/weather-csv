package main

import (
	"context"
	"github.com/aki-yogiri/weather-csv/handler"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ApiServerEnv struct {
	Host string
	Port string
}

func main() {

	var apiServerEnv ApiServerEnv
	envconfig.Process("API", &apiServerEnv)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var storeServerEnv handler.StoreServerEnv
	envconfig.Process("STORE", &storeServerEnv)

	e.GET("/weather", handler.DownloadWeatherCSV(storeServerEnv))

	e.Start(apiServerEnv.Host + ":" + apiServerEnv.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
