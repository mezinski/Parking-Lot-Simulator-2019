package handlers

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"gopkg.in/echo.v3"
)

//H ...
type H map[string]interface{}

//GetTasks ...
func GetTasks(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "tasks")
	}
}

//PutTask ...
func PutTask(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusCreated, H{
			"created": 123,
		})
	}
}

//DeleteTask ...
func DeleteTask(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		return c.JSON(http.StatusOK, H{
			"deleted": id,
		})
	}
}
