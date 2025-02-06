package price_service

type PriceData struct {
	Price       string `json:"price"`
	Conf        string `json:"conf"`
	Expo        int    `json:"expo"`
	PublishTime int    `json:"publish_time"`
}

type ParsedData struct {
	ID       string    `json:"id"`
	Price    PriceData `json:"price"`
	EMAPrice PriceData `json:"ema_price"`
	Metadata struct {
		Slot               int `json:"slot"`
		ProofAvailableTime int `json:"proof_available_time"`
		PrevPublishTime    int `json:"prev_publish_time"`
	} `json:"metadata"`
}

type ResponseData struct {
	Binary struct {
		Encoding string   `json:"encoding"`
		Data     []string `json:"data"`
	} `json:"binary"`
	Parsed []ParsedData `json:"parsed"`
}
