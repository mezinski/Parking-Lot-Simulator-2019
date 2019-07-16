package main

import (
	"Golang-Code/Go-with-Vue-2/handlers"
	"Golang-Code/Go-with-Vue-2/models"
	"fmt"

	"gopkg.in/echo.v3"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func main() {

	var todo handlers.Todos
	var vehicle handlers.Vehicles

	v, _ := InitConfig()

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", v.GetString("config.database.host"), v.GetString("config.database.port"), v.GetString("config.database.user"), v.GetString("config.database.dbname")))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.ParkedVehicle{})

	defer db.Close()

	todo.Db = db
	vehicle.Db = db
	vehicle.Cfg = v

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tasks", todo.GetTasks)
	e.PUT("tasks", todo.PutTask)
	e.DELETE("/tasks/:id", todo.DeleteTask)

	e.POST("tickets", vehicle.PostVehicleEntry)
	e.POST("/tickets/:id/:duration", vehicle.PostVehicleDuration)
	//e.GET("/tickets/:id")
	e.POST("/payments/:id", vehicle.PostVehiclePayment)

	e.Logger.Fatal(e.Start(":8000"))
}
