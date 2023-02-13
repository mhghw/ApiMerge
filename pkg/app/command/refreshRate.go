package command

import (
	"context"
	"mergeApi/config"
	"mergeApi/pkg/domain/rate"

	"github.com/sirupsen/logrus"
)

type RefreshRateHandler struct {
	provider map[string]RateGetter
	rateRepo rate.Repository
}

func NewRefreshRateHandler(
	ps map[string]RateGetter,
	rp rate.Repository,
) RefreshRateHandler {
	// if p == "" {
	// 	defer recover()
	// 	panic("invalid provider")
	// }
	// if rg == nil {
	// 	defer recover()
	// 	panic("nil rateGetter")
	// }
	handler := RefreshRateHandler{
		rateRepo: rp,
		provider: make(map[string]RateGetter),
	}
	if handler.rateRepo == nil {
		defer recover()
		panic("nil rate repo")
	}
	for k, v := range ps {
		handler.provider[k] = v
	}

	return handler
}

func (h RefreshRateHandler) Handle(ctx context.Context) {
	for k := range h.provider {
		updates, err := h.provider[k].GetUpdatedRates(ctx, config.AppConfig.Currencies)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		for _, update := range updates {
			err := h.rateRepo.RefreshRates(ctx, update)
			if err != nil {
				logrus.Errorln(err)
			}
		}
	}
}
