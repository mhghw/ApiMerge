package rate

import (
	"errors"
	"strings"
	"time"

	"mergeApi/config"

	"github.com/thoas/go-funk"
)

var (
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidProvider = errors.New("invalid provider")
)

// type Rate struct {
// 	from      string
// 	to        string
// 	exRate    float64
// 	provider  string
// 	createdAt time.Time
// }

type Rate struct {
	currency  string
	exRate    float64
	provider  string
	createdAt time.Time
}

func (r Rate) Currency() string {
	return r.currency
}
func (r Rate) ExchangeRate() float64 {
	return r.exRate
}
func (r Rate) Provider() string {
	return r.provider
}
func (r Rate) CreatedAt() time.Time {
	return r.createdAt
}

func NewRate(currency, provider string, exRate float64) (*Rate, error) {
	cur := strings.ToUpper(currency)
	if !funk.ContainsString(config.AppConfig.Currencies, cur) {
		return nil, ErrInvalidCurrency
	}

	return &Rate{
		currency:  cur,
		exRate:    exRate,
		provider:  provider,
		createdAt: time.Now().UTC(),
	}, nil
}

func UnmarshalFromDatabase(
	currency string,
	exRate float64,
	provider string,
	createdAt time.Time,
) *Rate {
	return &Rate{
		currency:  currency,
		exRate:    exRate,
		provider:  provider,
		createdAt: createdAt,
	}
}
