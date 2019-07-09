package main

import (
	"database/sql"
	"fmt"

	"gopkg.in/echo.v3"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// db := initDB("storage.db")
	// migrate(db)

	db, err := gorm.Open("postgres", "host=0.0.0.0 port=5432 user=postgres dbname=postgres")
	if err != nil {
		panic(err)
	}

	fmt.Println(err)
	defer db.Close()

	e := echo.New()

	e.GET("tasks", func(c echo.Context) error { return c.JSON(200, "GET Tasks") })
	e.PUT("tasks", func(c echo.Context) error { return c.JSON(200, "PUT Tasks") })
	e.DELETE("/tasks/:id", func(c echo.Context) error { return c.JSON(200, "DELETE Task "+c.Param("id")) })

	e.Logger.Fatal(e.Start(":8000"))
}

func initDB(filepath string) *sql.DB {

	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}

	return db
}

func migrate(db *sql.DB) {
	sql := `
		CREATE TABLE IF NOT EXISTS tasks(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name VARCHAR NOT NULL
		);
	`

	_, err := db.Exec(sql)

	if err != nil {
		panic(err)
	}
}
