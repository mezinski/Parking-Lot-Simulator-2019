package handlers

import (
	"Golang-Code/Go-with-Vue-2/models"
	"fmt"
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
		return c.JSON(http.StatusOK, models.GetTasks(db))
	}
}

//PutTask ...
func PutTask(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var task models.Task

		c.Bind(&task)

		id, err := models.PutTask(db, task.Name)

		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"created": id,
			})
		}
		return err
	}
}

//DeleteTask ...
func DeleteTask(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		_, err := models.DeleteTask(db, id)

		if err == nil {
			return c.JSON(http.StatusOK, H{
				"deleted": id,
			})
		}
		fmt.Println(err)
		return c.JSON(http.StatusOK, H{
			"error": fmt.Sprintf("%s", err),
		})
	}
}
