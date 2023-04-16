package api

import (
	"context"
	"time"
)

type FuelMessage struct {
	RecordID        string    `json:"recordid" gorm:"primaryKey"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Region          string    `json:"region"`
	Department      string    `json:"department"`
	DepCode         string    `json:"depcode"`
	City            string    `json:"city"`
	ZipCode         string    `json:"zipcode"`
	Address         string    `json:"address"`
	Name            string    `json:"name"`
	Price           float64   `json:"price"`
	PriceUpdateTime time.Time `json:"priceupdatetime"`
}

type Middleware interface {
	Send(Message FuelMessage) error
	Configure(ctx context.Context, c chan FuelMessage) error
}
