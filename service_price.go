package main

import "context"

type ServicePrice interface {
	MustQuerySuccess(int) float64
	GetPrice() float64
	Run(ctx context.Context)
}
