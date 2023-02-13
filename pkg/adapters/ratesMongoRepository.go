package adapters

import (
	"context"
	"log"
	"mergeApi/config"
	"mergeApi/pkg/domain/rate"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ratesMongoRepository struct {
	client *mongo.Client
}

func NewRatesMongoRepository(client *mongo.Client) ratesMongoRepository {
	repo := ratesMongoRepository{
		client: client,
	}

	model := RateModel{
		ExRate:    1,
		Currency:  config.AppConfig.BaseCurrency,
		Provider:  "Default",
		CreatedAt: time.Now().UTC(),
	}
	_, err := repo.RatesCollection().ReplaceOne(
		context.Background(),
		bson.M{"currency": model.Currency},
		model,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func (r ratesMongoRepository) marshalRate(rt *rate.Rate) RateModel {
	return RateModel{
		ExRate:    rt.ExchangeRate(),
		Currency:  rt.Currency(),
		Provider:  rt.Provider(),
		CreatedAt: rt.CreatedAt(),
	}
}
func (r ratesMongoRepository) unmarshalRate(rtm RateModel) *rate.Rate {
	return rate.UnmarshalFromDatabase(rtm.Currency, rtm.ExRate, rtm.Provider, rtm.CreatedAt)
}

type RateModel struct {
	ExRate    float64   `bson:"exRate"`
	Currency  string    `bson:"currency"`
	Provider  string    `bson:"provider"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (r ratesMongoRepository) RatesCollection() *mongo.Collection {
	return r.client.Database(viper.GetString("mongo.db")).Collection("rates")
}

func (r ratesMongoRepository) RefreshRates(ctx context.Context, rt *rate.Rate) error {
	model := r.marshalRate(rt)
	filter := bson.M{"$and": bson.A{
		bson.D{{"currency", strings.ToUpper(model.Currency)}},
		bson.D{{"provider", model.Provider}},
	}}
	opts := options.Replace().SetUpsert(true)

	_, err := r.RatesCollection().ReplaceOne(ctx, filter, model, opts)
	return err

}
func (r ratesMongoRepository) GetLowestRate(ctx context.Context, currency string) (*rate.Rate, error) {
	filter := bson.M{"currency": strings.ToUpper(currency)}

	cur, err := r.RatesCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	rates := make([]RateModel, 0)
	cur.All(ctx, &rates)

	return r.unmarshalRate(r.findLowest(rates)), nil
}

func (r ratesMongoRepository) findLowest(rates []RateModel) RateModel {
	if len(rates) == 0 {
		return RateModel{}
	}
	lowest := rates[0]
	for i := 0; i < len(rates); i++ {
		if lowest.ExRate < rates[i].ExRate {
			continue
		} else {
			lowest = rates[i]
		}
	}

	return lowest
}
