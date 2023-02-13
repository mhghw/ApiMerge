package adapters

import (
	"context"
	"encoding/json"
	"mergeApi/config"
	"mergeApi/pkg/domain/rate"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type GeckoHttpClient struct {
	url      string
	provider string
}

func (c GeckoHttpClient) Provider() string {
	return c.provider
}

func NewGeckoHttpClient() GeckoHttpClient {
	return GeckoHttpClient{
		url:      viper.GetString("provider.gecko"),
		provider: "Gecko",
	}
}

func (c GeckoHttpClient) GetUpdatedRates(ctx context.Context, currencies []string) ([]*rate.Rate, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body geckoResponseBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	usdRate := body.Rates["usd"].(map[string]any)["value"].(float64)
	supportedRates := make([]*rate.Rate, 0)
	for _, cur := range config.AppConfig.Currencies {
		r := body.Rates[strings.ToLower(cur)]
		v := usdRate / r.(map[string]any)["value"].(float64)
		supRate, _ := rate.NewRate(cur, c.provider, v)
		supportedRates = append(supportedRates, supRate)
	}

	return supportedRates, nil

}

type geckoResponseBody struct {
	Rates map[string]any `json:"rates"`
}
