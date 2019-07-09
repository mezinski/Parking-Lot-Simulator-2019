package main

import (
	"fmt"

	"Golang-Code/Go-with-Vue-2/models"

	"gopkg.in/echo.v3"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// db := initDB("storage.db")
	// migrate(db)

	db, err := gorm.Open("postgres", "host=0.0.0.0 port=5439 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Task{})

	fmt.Println(err)
	defer db.Close()

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("tasks", handlers.GetTasks(db))
	e.PUT("tasks", handlers.PutTask(db))
	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	e.Logger.Fatal(e.Start(":8000"))
}

// func initDB(filepath string) *sql.DB {

// 	db, err := sql.Open("sqlite3", filepath)

// 	if err != nil {
// 		panic(err)
// 	}

// 	if db == nil {
// 		panic("db nil")
// 	}

// 	return db
// }

// func migrate(db *sql.DB) {
// 	sql := `
// 		CREATE TABLE IF NOT EXISTS tasks(
// 			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
// 			name VARCHAR NOT NULL
// 		);
// 	`

// 	_, err := db.Exec(sql)

// 	if err != nil {
// 		panic(err)
// 	}
// }
