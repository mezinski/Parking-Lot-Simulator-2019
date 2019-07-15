package handlers

import (
	"Golang-Code/Go-with-Vue-2/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"gopkg.in/echo.v3"
)

//Todos ...
type Todos struct {
	Db *gorm.DB
}

//H ...
type H map[string]interface{}

//GetTasks ...
func (t *Todos) GetTasks(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetTasks(t.Db))
}

//PutTask ...
func (t *Todos) PutTask(c echo.Context) error {

	var task models.Task

	c.Bind(&task)

	id, err := models.PutTask(t.Db, task.Name)

	if err == nil {
		return c.JSON(http.StatusCreated, H{
			"created": id,
		})
	}
	return err
}

//DeleteTask ...
func (t *Todos) DeleteTask(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	_, err := models.DeleteTask(t.Db, id)

	if err == nil {
		return c.JSON(http.StatusOK, H{
			"delete_id": id,
			"response":  fmt.Sprintf("%s", err),
		})
	}
	return c.JSON(http.StatusOK, H{
		"delete_id": id,
		"response":  fmt.Sprintf("%s", err),
	})
}
