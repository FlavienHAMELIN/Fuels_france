package adapters

import (
	"context"
	"fmt"

	"github.com/FlavienHAMELIN/Fuels_france/backend/api"
	"github.com/FlavienHAMELIN/Fuels_france/backend/config"
	"github.com/go-logr/logr"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseAdapter struct {
	Logger         logr.Logger
	Db             *gorm.DB
	ReceiveChannel chan api.FuelMessage
}

func (g *DatabaseAdapter) Send(message api.FuelMessage) error {
	g.Logger.Info(
		"Données",
		"RecordID", message.RecordID,
		"Name", message.Name,
		"Price", message.Price,
		"Latitude", message.Latitude,
		"Longitude", message.Longitude,
		"Region", message.Region,
		"Departement", message.Department,
		"DepCode", message.DepCode,
		"Address", message.Address,
		"City", message.City,
		"ZipCode", message.ZipCode,
	)

	g.Db.Save(&message)

	return nil
}

func (d *DatabaseAdapter) Configure(ctx context.Context, c chan api.FuelMessage) error {
	d.ReceiveChannel = c

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Configuration.Database.Host,
		config.Configuration.Database.Port,
		config.Configuration.Database.User,
		config.Configuration.Database.Password,
		config.Configuration.Database.Db,
	)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	d.Db = db
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&api.FuelMessage{})

	go d.InsertUpdateAuto()

	return nil
}

func (d DatabaseAdapter) InsertUpdateAuto() {
	for message := range d.ReceiveChannel {
		d.Send(message)
	}
}

func (dbA DatabaseAdapter) GetData(args map[string]string) []api.FuelMessage {
	var result []api.FuelMessage
	dbA.Db.Where(args).Find(&result)
	return result
}

func (dbA DatabaseAdapter) GetDataLowerPrice(args map[string]string) []api.FuelMessage {
	var result []api.FuelMessage
	fmt.Println(args)
	dbA.Db.Where(args).Order("price ASC").Find(&result)
	return result
}

func (dbA DatabaseAdapter) GetDataLowerPriceLimit(args map[string]string) []api.FuelMessage {
	var result []api.FuelMessage
	fmt.Println(args)
	dbA.Db.Where(args).Order("price ASC").Limit(20).Find(&result)
	return result
}
