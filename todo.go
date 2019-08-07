package main

import (
	"gopkg.in/echo.v3"

	_ "github.com/lib/pq"
	"github.com/mezinski/Go-with-Vue-2/handlers"
	"github.com/mezinski/Go-with-Vue-2/models"
)

func main() {

	var vehicle handlers.Vehicles

	v, _ := InitConfig()
	db, _ := DbConnection(v)

	db.AutoMigrate(&models.ParkedVehicle{})

	defer db.Close()

	vehicle.Db = db
	vehicle.Cfg = v

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tickets", vehicle.GetVehicles)
	e.GET("/tickets/:id", vehicle.GetVehicleByID)
	e.POST("tickets", vehicle.PostVehicleEntry)
	e.POST("/payments/:id", vehicle.PostVehiclePayment)

	e.Logger.Fatal(e.Start(":8000"))
}
