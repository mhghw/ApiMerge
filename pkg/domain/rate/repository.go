package rate

import "context"

type Repository interface {
	RefreshRates(ctx context.Context, rt *Rate) error
	GetLowestRate(ctx context.Context, currency string) (*Rate, error)
}
