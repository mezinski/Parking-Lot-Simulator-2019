package main

import (
	"gopkg.in/echo.v3"

	_ "github.com/lib/pq"
	"github.com/mezinski/Parking-Lot-Simulator-2019/models"
)

func main() {

	var routes Routing

	v, _ := InitConfig()
	db, _ := DbConnection(v)

	db.AutoMigrate(&models.Vehicle{})

	defer db.Close()

	routes.Db = db
	routes.Cfg = v

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tickets", routes.GetVehicles)
	e.GET("/tickets/:id", routes.GetVehicleByID)
	e.POST("tickets", routes.CreateVehicle)
	e.POST("/payments/:id", routes.PostVehiclePayment)

	e.Logger.Fatal(e.Start(":8000"))
}
