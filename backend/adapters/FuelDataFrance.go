package adapters

import (
	"encoding/json"
	"flag"
	"strconv"
	"strings"
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
	ID        int    `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Cp        string `json:"cp"`
	Pop       string `json:"pop"`
	Adresse   string `json:"adresse"`
	Ville     string `json:"ville"`
	Horaires  string `json:"horaires"`
	Services  string `json:"services"`
	Prix      string `json:"prix"`
	Rupture   string `json:"rupture"`
	Geom      struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"geom"`
	GazoleMaj               time.Time `json:"gazole_maj"`
	GazolePrix              float64   `json:"gazole_prix"`
	Sp95Maj                 time.Time `json:"sp95_maj"`
	Sp95Prix                float64   `json:"sp95_prix"`
	E85Maj                  time.Time `json:"e85_maj"`
	E85Prix                 float64   `json:"e85_prix"`
	GplcMaj                 time.Time `json:"gplc_maj"`
	GplcPrix                float64   `json:"gplc_prix"`
	E10Maj                  time.Time `json:"e10_maj"`
	E10Prix                 float64   `json:"e10_prix"`
	Sp98Maj                 time.Time `json:"sp98_maj"`
	Sp98Prix                float64   `json:"sp98_prix"`
	CarburantsDisponibles   []string  `json:"carburants_disponibles"`
	CarburantsIndisponibles []string  `json:"carburants_indisponibles"`
	Departement             string    `json:"departement"`
	CodeDepartement         string    `json:"code_departement"`
	Region                  string    `json:"region"`
	CodeRegion              string    `json:"code_region"`
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
		if len(item.CarburantsDisponibles) != 0 {

			for i, fuel_name := range item.CarburantsDisponibles {
				var recordid string
				var region string
				var departement string
				var depcode string
				var city string
				var zipcode string
				var address string
				var name string
				var price float64
				var priceupdatetime time.Time

				recordid = strconv.Itoa(item.ID + i)
				region = item.Region
				departement = item.Departement
				depcode = item.CodeDepartement
				city = item.Ville
				zipcode = item.Cp
				address = item.Adresse
				name = fuel_name
				if strings.Contains(strings.ToLower(name), "gazole") {
					price = item.GazolePrix
					priceupdatetime = item.GazoleMaj
				}
				if strings.Contains(strings.ToLower(name), "sp95") {
					price = item.Sp95Prix
					priceupdatetime = item.Sp95Maj
				}
				if strings.Contains(strings.ToLower(name), "e85") {
					price = item.E85Prix
					priceupdatetime = item.E85Maj
				}
				if strings.Contains(strings.ToLower(name), "gplc") {
					price = item.GplcPrix
					priceupdatetime = item.GplcMaj
				}
				if strings.Contains(strings.ToLower(name), "e10") {
					price = item.E10Prix
					priceupdatetime = item.E10Maj
				}
				if strings.Contains(strings.ToLower(name), "sp98") {
					price = item.Sp98Prix
					priceupdatetime = item.Sp98Maj
				}
				if priceupdatetime.After(before) {
					fuel := api.FuelMessage{
						RecordID:        recordid,
						Region:          region,
						Department:      departement,
						DepCode:         depcode,
						City:            city,
						ZipCode:         zipcode,
						Address:         address,
						Name:            name,
						Price:           price,
						PriceUpdateTime: priceupdatetime,
					}

					if item.Geom.Lat != 0 && item.Geom.Lon != 0 {
						fuel.Latitude = item.Geom.Lat
						fuel.Longitude = item.Geom.Lon
					}

					(*a.middleware).Send(fuel)
				}
			}
		}

	}
	flag.Set("logtostderr", "true")
	a.Logger.Info("Database updated")
	flag.Set("logtostderr", "false")

	return nil
}
