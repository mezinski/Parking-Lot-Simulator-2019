package handlers

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/mezinski/Go-with-Vue-2/models"
	"github.com/spf13/viper"
	"gopkg.in/echo.v3"
)

//Vehicles - This is a struct for vehicles used for handling. We attach the DB object and the Config object to it so we can utilize these intuitively from within the methods they're attributed to.
type Vehicles struct {
	Db  *gorm.DB
	Cfg *viper.Viper
}

//H This is used to pass back JSON blobs to the requestor
type H map[string]interface{}

//GetVehicles - Handler method used to get all vehicles currently parked in the DB
func (v *Vehicles) GetVehicles(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetVehicles(v.Db))
}

//GetVehicleByID - Handler method used to get a vehicle parked in the DB, by ID
func (v *Vehicles) GetVehicleByID(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	licensePlate, duration, totalPaid := models.GetVehicleByID(v.Db, id)

	return c.JSON(http.StatusOK, H{
		"id":            id,
		"license_plate": licensePlate,
		"duration":      duration,
		"total_paid":    totalPaid,
	})
}

//PostVehicleEntry - Handler method to bind to a Vehicle object, create the vehicle object, and return a response to the requestor
func (v *Vehicles) PostVehicleEntry(c echo.Context) error {

	var vehicle models.ParkedVehicle

	c.Bind(&vehicle)

	newVehicle, err := models.PostVehicleEntry(v.Db, v.Cfg, vehicle.LicensePlate)
	if err == nil {
		return c.JSON(http.StatusCreated, newVehicle)
	}
	return c.JSON(http.StatusOK, H{
		"created": 0,
		"vehicle": newVehicle,
		"error":   err.Error(),
	})
}

//PostVehiclePayment - Handler method used to take in a request for payment, 'process' the payment, and 'remove' the car from the lot
func (v *Vehicles) PostVehiclePayment(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	_, licensePlate, duration, totalPaid, err := models.PostVehiclePayment(v.Db, v.Cfg, id)

	if err == nil {
		return c.JSON(http.StatusOK, H{
			"id":            id,
			"license_plate": licensePlate,
			"duration":      duration,
			"total_paid":    totalPaid,
		})
	}
	return c.JSON(http.StatusOK, H{
		"message": err.Error(),
	})
}
