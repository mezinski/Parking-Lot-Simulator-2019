package main

import (
	"gopkg.in/echo.v3"

	_ "github.com/lib/pq"
	"github.com/mezinski/Go-with-Vue-2/models"
)

func main() {

	var route Routing

	v, _ := InitConfig()
	db, _ := DbConnection(v)

	db.AutoMigrate(&models.Vehicle{})

	defer db.Close()

	route.Db = db
	route.Cfg = v

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tickets", route.GetVehicles)
	e.GET("/tickets/:id", route.GetVehicleByID)
	e.POST("tickets", route.PostVehicleEntry)
	e.POST("/payments/:id", route.PostVehiclePayment)

	e.Logger.Fatal(e.Start(":8000"))
}
