package query

import (
	"context"
	"mergeApi/pkg/domain/rate"

	"github.com/sirupsen/logrus"
)

type GetRatesHandler struct {
	readModel GetRatesReadModel
}

func NewGetRatesHandler(readModel GetRatesReadModel) GetRatesHandler {
	if readModel == nil {
		defer recover()
		panic("nil readModel")
	}

	return GetRatesHandler{
		readModel: readModel,
	}
}

type GetRatesReadModel interface {
	GetLowestRate(ctx context.Context, currency string) (*rate.Rate, error)
}

func (h GetRatesHandler) Handle(ctx context.Context, currencies []string, to string) ([]ExchangeRate, error) {
	target, err := h.readModel.GetLowestRate(ctx, to)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	exRates := make([]ExchangeRate, 0)
	for _, src := range currencies {
		lowest, err := h.readModel.GetLowestRate(ctx, src)
		if err != nil {
			logrus.Errorln(err)
			return nil, err
		}
		exRate := ExchangeRate{
			To:       target.Currency(),
			From:     src,
			Rate:     lowest.ExchangeRate() / target.ExchangeRate(),
			Provider: lowest.Provider(),
		}
		exRates = append(exRates, exRate)
	}

	return exRates, nil
}
