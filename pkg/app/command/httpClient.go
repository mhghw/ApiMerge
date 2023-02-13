package command

import (
	"context"
	"mergeApi/pkg/domain/rate"
)

type RateGetter interface {
	GetUpdatedRates(ctx context.Context, currencies []string) ([]*rate.Rate, error)
}
