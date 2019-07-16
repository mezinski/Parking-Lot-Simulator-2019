package handlers

import (
	"Golang-Code/Go-with-Vue-2/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"gopkg.in/echo.v3"
)

//Todos ...
type Todos struct {
	Db *gorm.DB
}

//Vehicles ...
type Vehicles struct {
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

//PostVehicleEntry ...
func (v *Vehicles) PostVehicleEntry(c echo.Context) error {

	var vehicle models.ParkedVehicle

	c.Bind(&vehicle)

	id, err := models.PostVehicleEntry(v.Db, vehicle.LicensePlate)
	fmt.Println(id)
	fmt.Println(err)
	if err == nil {
		return c.JSON(http.StatusCreated, H{
			"created":       id,
			"license_plate": vehicle.LicensePlate,
		})
	}
	return c.JSON(http.StatusOK, H{
		"error": err.Error(),
	})
}

//PostVehicleDuration ...
func (v *Vehicles) PostVehicleDuration(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	duration, _ := strconv.Atoi(c.Param("duration"))

	_, licensePlate, duration, err := models.PostVehicleDuration(v.Db, id, duration)

	if err == nil {
		message := fmt.Sprintf("Your ticket is valid for %d hours. Please pay on your way out.", duration)
		return c.JSON(http.StatusOK, H{
			"id":            id,
			"license_plate": licensePlate,
			"duration":      duration,
			"message":       message,
		})
	}
	return c.JSON(http.StatusOK, H{
		"message": err.Error(),
	})
}

//PostVehiclePayment ...
func (v *Vehicles) PostVehiclePayment(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	_, licensePlate, duration, totalPaid, err := models.PostVehiclePayment(v.Db, id)

	if err == nil {
		message := fmt.Sprintf("Thank you, %s, for choosing sketchypark. Payment has been processed by your %s, Card no. %d", viper.GetString("config.user-data.name-on-card"), viper.GetString("config.user-data.credit-card-type"), viper.GetInt("config.user-data.credit-card-no"))
		return c.JSON(http.StatusOK, H{
			"id":            id,
			"license_plate": licensePlate,
			"duration":      duration,
			"total_paid":    totalPaid,
			"message":       message,
		})
	}
	return c.JSON(http.StatusOK, H{
		"message": err.Error(),
	})
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
