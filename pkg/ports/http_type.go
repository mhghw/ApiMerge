package ports

type Response struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Provider string  `json:"provider"`
	Rate     float64 `json:"rate"`
}

type RequestBody struct {
	Currencies []string `json:"currencies"`
	To         string   `json:"to"`
}
