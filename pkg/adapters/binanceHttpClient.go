package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"mergeApi/pkg/domain/rate"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type BinanceHttpClient struct {
	url      string
	provider string
}

func (c BinanceHttpClient) Provider() string {
	return c.provider
}

func NewBinanceHttpClient() BinanceHttpClient {
	return BinanceHttpClient{
		url:      viper.GetString("provider.Binance"),
		provider: "Binance",
	}
}

func (c BinanceHttpClient) GetUpdatedRates(ctx context.Context, currencies []string) ([]*rate.Rate, error) {
	rates := make([]*rate.Rate, 0)
	for _, cur := range currencies {
		symbol := c.symbolMaker(cur)
		url := fmt.Sprintf("%v?symbol=%v", viper.GetString("provider.binance"), symbol)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var body binanceResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return nil, err
		}
		ex, err := strconv.ParseFloat(body.Price, 64)
		if err != nil {
			return nil, err
		}
		r, err := rate.NewRate(cur, c.provider, ex)
		if err != nil {
			return nil, err
		}
		rates = append(rates, r)
	}

	return rates, nil

}

func (c BinanceHttpClient) symbolMaker(currency string) string {
	return fmt.Sprintf("%vUSDT", strings.ToUpper(currency))
}

type binanceResponse struct {
	Price string `json:"price"`
}
