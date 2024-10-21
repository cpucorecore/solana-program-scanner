package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

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

type ServicePriceHermes struct {
	QueryInterval time.Duration
	QueryApiUrl   string
	SolUsdPrice   float64
	Mu            sync.Mutex
}

func (sph *ServicePriceHermes) MustQuerySuccess(maxRetry int) float64 {
	sph.MustQueryPrice(maxRetry)
	return sph.SolUsdPrice
}

func (sph *ServicePriceHermes) GetPrice() float64 {
	sph.Mu.Lock()
	defer sph.Mu.Unlock()
	return sph.SolUsdPrice
}

func (sph *ServicePriceHermes) Run(ctx context.Context) {
	ticker := time.NewTicker(sph.QueryInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				sph.queryPrice()
			}
		}
	}()
}

func (sph *ServicePriceHermes) MustQueryPrice(maxRetry int) {
	for i := 0; i < maxRetry; i++ {
		err := sph.queryPrice()
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}

	Logger.Fatal(fmt.Sprintf("query price failed after %d retries", maxRetry))
}

func calculatePrice(price string, expo int) (float64, error) {
	priceInt := 0.0
	_, err := fmt.Sscanf(price, "%f", &priceInt)
	if err != nil {
		return 0.0, err
	}

	calculatedPrice := priceInt * math.Pow(10, float64(expo))
	return calculatedPrice, nil
}

func (sph *ServicePriceHermes) queryPrice() error {
	startTime := time.Now()
	Logger.Info(fmt.Sprintf("Querying price begin at %v", startTime))
	resp, err := http.Get(sph.QueryApiUrl)
	endTime := time.Now()
	if err != nil {
		Logger.Error(fmt.Sprintf("Querying price err: %v at %v", err, endTime))
		return err
	}
	Logger.Info(fmt.Sprintf("Querying price success at %v, used %d ms", startTime, endTime.Sub(startTime).Milliseconds()))

	var data ResponseData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	if len(data.Parsed) == 0 {
		return fmt.Errorf("no parsed data found")
	}

	priceData := data.Parsed[0].Price
	expo := priceData.Expo
	price, err := calculatePrice(priceData.Price, expo)
	if err != nil {
		return err
	}

	sph.Mu.Lock()
	defer sph.Mu.Unlock()
	sph.SolUsdPrice = price
	resp.Body.Close()

	return nil
}

var _ ServicePrice = &ServicePriceHermes{}

func NewServicePriceHermes(queryInterval time.Duration, queryApiUrl string) *ServicePriceHermes {
	return &ServicePriceHermes{
		QueryInterval: queryInterval,
		QueryApiUrl:   queryApiUrl,
	}
}
