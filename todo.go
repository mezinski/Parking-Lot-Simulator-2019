package main

import (
	"Golang-Code/Go-with-Vue-2/handlers"
	"Golang-Code/Go-with-Vue-2/models"

	"gopkg.in/echo.v3"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := gorm.Open("postgres", "host=0.0.0.0 port=5439 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Task{})

	defer db.Close()

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tasks", handlers.GetTasks(db))
	e.PUT("tasks", handlers.PutTask(db))
	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	e.Logger.Fatal(e.Start(":8000"))
}
