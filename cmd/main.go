package main

import (
	"context"
	"log"
	"mergeApi/config"
	"mergeApi/pkg/ports"
	"mergeApi/pkg/service"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.InitConfig()

	ticker := time.NewTicker(time.Duration(config.AppConfig.Interval * 1000000000))
	app := service.NewApplication(context.Background())

	srv := ports.NewHttpServer(app)

	srv.Run(ticker, ":8080")

}
