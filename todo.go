package main

import (
	"fmt"

	"gopkg.in/echo.v3"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/mezinski/Go-with-Vue-2/handlers"
	"github.com/mezinski/Go-with-Vue-2/models"
)

func main() {

	var todo handlers.Todos
	var vehicle handlers.Vehicles

	v, _ := InitConfig()

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", v.GetString("config.database.host"), v.GetString("config.database.port"), v.GetString("config.database.user"), v.GetString("config.database.dbname")))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.ParkedVehicle{})

	defer db.Close()

	todo.Db = db
	vehicle.Db = db
	vehicle.Cfg = v

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tickets", vehicle.GetVehicles)
	e.GET("/tickets/:id", vehicle.GetVehicleByID)
	e.POST("tickets", vehicle.PostVehicleEntry)
	e.POST("/tickets/:id/:duration", vehicle.PostVehicleDuration)
	//e.GET("/tickets/:id")
	e.POST("/payments/:id", vehicle.PostVehiclePayment)

	e.Logger.Fatal(e.Start(":8000"))
}
