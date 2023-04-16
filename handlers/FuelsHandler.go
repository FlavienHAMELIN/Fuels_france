package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/FlavienHAMELIN/fuelPricesFrance/adapters"
)

func FuelsHandler(dbA *adapters.DatabaseAdapter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		args := make(map[string]string)
		record_id := r.URL.Query().Get("record_id")
		if record_id != "" {
			args["record_id"] = record_id
		}
		region := r.URL.Query().Get("region")
		if region != "" {
			args["region"] = region
		}
		department := r.URL.Query().Get("department")
		if department != "" {
			args["department"] = department
		}
		dep_code := r.URL.Query().Get("dep_code")
		if dep_code != "" {
			args["dep_code"] = dep_code
		}
		city := r.URL.Query().Get("city")
		if city != "" {
			args["city"] = city
		}
		zip_code := r.URL.Query().Get("zip_code")
		if zip_code != "" {
			args["zip_code"] = zip_code
		}
		address := r.URL.Query().Get("address")
		if address != "" {
			args["address"] = address
		}
		name := r.URL.Query().Get("name")
		if name != "" {
			args["name"] = name
		}
		lower_price := r.URL.Query().Get("lower_price")
		limit := r.URL.Query().Get("limit")

		if len(args) != 0 {
			if lower_price == "1" {
				if limit == "1" {
					datas := dbA.GetDataLowerPriceLimit(args)
					result, err := json.Marshal(datas)
					if err != nil {
						panic(err)
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(result)
				} else {
					datas := dbA.GetDataLowerPrice(args)
					result, err := json.Marshal(datas)
					if err != nil {
						panic(err)
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(result)
				}

			} else {
				datas := dbA.GetData(args)
				result, err := json.Marshal(datas)
				if err != nil {
					panic(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(result)
			}
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
