package service

import (
	"context"
	"log"
	"mergeApi/pkg/adapters"
	"mergeApi/pkg/app"
	"mergeApi/pkg/app/command"
	"mergeApi/pkg/app/query"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewApplication(ctx context.Context) app.Application {
	// storage := adapters.NewRatesMemoryRepository()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatal(err)
	}
	ratesRepo := adapters.NewRatesMongoRepository(mongoClient)
	geckoClient := adapters.NewGeckoHttpClient()
	binanceClient := adapters.NewBinanceHttpClient()
	providerMap := map[string]command.RateGetter{
		geckoClient.Provider():   geckoClient,
		binanceClient.Provider(): binanceClient,
	}

	return app.Application{
		Commands: app.Command{
			RefreshRate: command.NewRefreshRateHandler(providerMap, ratesRepo),
		},

		Queries: app.Query{
			GetRates: query.NewGetRatesHandler(ratesRepo),
		},
	}
}
