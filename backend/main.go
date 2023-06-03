package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/FlavienHAMELIN/Fuels_france/backend/adapters"
	"github.com/FlavienHAMELIN/Fuels_france/backend/api"
	"github.com/FlavienHAMELIN/Fuels_france/backend/handlers"
	"github.com/go-logr/glogr"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

var Logger logr.Logger

func main() {
	readFuel := make(chan api.FuelMessage)

	databaseAdapter := &adapters.DatabaseAdapter{
		Logger: Logger.WithName("databaseMiddleware"),
	}
	var databaseMiddleware api.Middleware = databaseAdapter
	databaseMiddleware.Configure(context.Background(), readFuel)

	t := time.NewTicker(1 * time.Second)
	go func() {
		for ; true; <-t.C {
			hour, min, sec := time.Now().Clock()
			if hour == 12 && min == 00 && sec == 00 {
				var api api.FuelsData = adapters.APIFuel{
					Logger: Logger.WithName("APIfuels"),
				}
				api.GetFuels(databaseMiddleware)
			}
		}
	}()

	var api api.FuelsData = adapters.APIFuel{
		Logger: Logger.WithName("APIfuels"),
	}
	api.GetFuels(databaseMiddleware)

	r := mux.NewRouter()
	r.HandleFunc("/fuelsFrance", handlers.FuelsHandler(databaseAdapter))
	r.Headers()
	channel := make(chan struct{})
	go http.ListenAndServe(":8181", r)
	flag.Set("logtostderr", "true")
	Logger.Info("Listening at localhost:8181/fuelsFrance")
	<-channel
}

func init() {
	flag.Parse()
	flag.Set("log_dir", "logs/")
	Logger = glogr.NewWithOptions(glogr.Options{}).WithName("fuelsFrance")
}
