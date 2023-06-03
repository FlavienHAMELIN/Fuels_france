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
	Datasetid string `json:"datasetid"`
	Recordid  string `json:"recordid"`
	Fields    struct {
		Geom                 []float64 `json:"geom"`
		RegName              string    `json:"reg_name"`
		PrixValeur           float64   `json:"prix_valeur"`
		ComCode              string    `json:"com_code"`
		DepName              string    `json:"dep_name"`
		Horaires             string    `json:"horaires"`
		Adresse              string    `json:"adresse"`
		PrixNom              string    `json:"prix_nom"`
		ComName              string    `json:"com_name"`
		Ville                string    `json:"ville"`
		ComArmName           string    `json:"com_arm_name"`
		RegCode              string    `json:"reg_code"`
		ServicesService      string    `json:"services_service"`
		DepCode              string    `json:"dep_code"`
		ComArmCode           string    `json:"com_arm_code"`
		EpciCode             string    `json:"epci_code"`
		Cp                   string    `json:"cp"`
		ID                   string    `json:"id"`
		PrixID               string    `json:"prix_id"`
		EpciName             string    `json:"epci_name"`
		Pop                  string    `json:"pop"`
		PrixMaj              time.Time `json:"prix_maj"`
		HorairesAutomate2424 string    `json:"horaires_automate_24_24"`
	} `json:"fields"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
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
		if item.Fields.PrixValeur != 0 && item.Fields.PrixMaj.After(before) {
			fuel := api.FuelMessage{
				RecordID:        item.Recordid,
				Region:          item.Fields.RegName,
				Department:      item.Fields.DepName,
				DepCode:         item.Fields.DepCode,
				City:            item.Fields.Ville,
				ZipCode:         item.Fields.Cp,
				Address:         item.Fields.Adresse,
				Name:            item.Fields.PrixNom,
				Price:           item.Fields.PrixValeur,
				PriceUpdateTime: item.Fields.PrixMaj,
			}

			if len(item.Fields.Geom) == 2 {
				fuel.Latitude = item.Fields.Geom[0]
				fuel.Longitude = item.Fields.Geom[1]
			}

			(*a.middleware).Send(fuel)
		}

	}
	flag.Set("logtostderr", "true")
	a.Logger.Info("Database updated")
	flag.Set("logtostderr", "false")

	return nil
}
