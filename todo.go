package main

import (
	"Golang-Code/Go-with-Vue-2/handlers"
	"Golang-Code/Go-with-Vue-2/models"

	"gopkg.in/echo.v3"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func main() {

	var todo handlers.Todos

	db, err := gorm.Open("postgres", "host=0.0.0.0 port=5439 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Task{})

	defer db.Close()

	todo.Db = db

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tasks", todo.GetTasks)
	e.PUT("tasks", todo.PutTask)
	e.DELETE("/tasks/:id", todo.DeleteTask)

	e.Logger.Fatal(e.Start(":8000"))
}
