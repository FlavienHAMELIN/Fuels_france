package adapters

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/FlavienHAMELIN/Fuels_france/backend/api"
	"github.com/FlavienHAMELIN/Fuels_france/backend/config"
	"github.com/go-logr/logr"
	"github.com/vicanso/go-axios"
)

type APIFuel struct {
	Logger     logr.Logger
	middleware *api.Middleware
}

type Fuels []struct {
	ID       string `json:"id"`
	Cp       string `json:"cp"`
	Pop      string `json:"pop"`
	Adresse  string `json:"adresse"`
	Ville    string `json:"ville"`
	Horaires string `json:"horaires"`
	Geom     struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"geom"`
	PrixMaj              time.Time `json:"prix_maj"`
	PrixID               string    `json:"prix_id"`
	PrixValeur           float64   `json:"prix_valeur"`
	PrixNom              string    `json:"prix_nom"`
	ComArmCode           string    `json:"com_arm_code"`
	ComArmName           string    `json:"com_arm_name"`
	EpciCode             string    `json:"epci_code"`
	EpciName             string    `json:"epci_name"`
	DepCode              string    `json:"dep_code"`
	DepName              string    `json:"dep_name"`
	RegCode              string    `json:"reg_code"`
	RegName              string    `json:"reg_name"`
	ComCode              string    `json:"com_code"`
	ComName              string    `json:"com_name"`
	ServicesService      []string  `json:"services_service"`
	HorairesAutomate2424 string    `json:"horaires_automate_24_24"`
}

func (a APIFuel) GetFuels(middleware api.Middleware) error {
	a.middleware = &middleware

	flag.Set("logtostderr", "true")
	a.Logger.Info("Connecting to the API", "URL", config.Configuration.API.URL)
	resp, err := axios.Get(config.Configuration.API.URL)
	if err != nil {
		panic(err)
	}
	a.Logger.Info("Successful connection to the API")

	var responseObject Fuels
	err = json.Unmarshal(resp.Data, &responseObject)
	if err != nil {
		panic(err)
	}

	a.Logger.Info("Updating the database")
	flag.Set("logtostderr", "false")
	today := time.Now()
	before := today.Add(-72 * time.Hour)
	for _, item := range responseObject {
		if item.PrixValeur != 0 && item.PrixMaj.After(before) {
			fuel := api.FuelMessage{
				RecordID:        item.ID + item.PrixID,
				Region:          item.RegName,
				Department:      item.DepName,
				DepCode:         item.DepCode,
				City:            item.Ville,
				ZipCode:         item.Cp,
				Address:         item.Adresse,
				Name:            item.PrixNom,
				Price:           item.PrixValeur,
				PriceUpdateTime: item.PrixMaj,
			}

			if item.Geom.Lat != 0 && item.Geom.Lon != 0 {
				fuel.Latitude = item.Geom.Lat
				fuel.Longitude = item.Geom.Lon
			}

			(*a.middleware).Send(fuel)
		}

	}
	flag.Set("logtostderr", "true")
	a.Logger.Info("Database updated")
	flag.Set("logtostderr", "false")

	return nil
}
