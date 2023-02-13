package ports

import (
	"context"
	"errors"
	"log"
	"mergeApi/pkg/app"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
)

type HttpServer struct {
	engine *gin.Engine
	app    app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	srv := HttpServer{
		engine: gin.Default(),
		app:    app,
	}

	srv.engine.GET("/getRate", srv.HandleGetRates)

	return srv
}

func (s HttpServer) Run(ticker *time.Ticker, addr ...string) error {
	var g run.Group
	g.Add(func() error {
		return s.engine.Run(":8000")
	}, func(err error) {
		log.Fatal(err)
	})
	g.Add(func() error {
		s.app.Commands.RefreshRate.Handle(context.Background())
		for range ticker.C {
			s.app.Commands.RefreshRate.Handle(context.Background())
		}
		return errors.New("worker stopped")
	}, func(err error) {
		log.Fatal(err)
	})

	return g.Run()
}
