package models

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

//ParkedVehicle - This is the struct for a 'Parked Vehicle'. This is all the information we store and utilize for a vehicle parking with this application.
type ParkedVehicle struct {
	ID           uint    `json:"id"`
	LicensePlate string  `json:"license_plate"`
	Duration     int     `json:"duration"`
	IsParked     bool    `json:"is_parked"`
	TotalPaid    float64 `json:"total_paid"`
}

//VehicleCollection - This struct is to hold an object which is an array of vehicle object types.
type VehicleCollection struct {
	ParkedVehicles []ParkedVehicle `json:"parked_vehicles"`
}

//GetVehicles - This method will return all of the vehicles currently parked in the lot
func GetVehicles(db *gorm.DB) VehicleCollection {

	rows, err := db.Table("parked_vehicles").Where("is_parked = true").Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	result := VehicleCollection{}

	for rows.Next() {
		vehicle := ParkedVehicle{}
		err2 := rows.Scan(&vehicle.ID, &vehicle.LicensePlate, &vehicle.Duration, &vehicle.IsParked, &vehicle.TotalPaid)
		if err2 != nil {
			panic(err2)
		}
		result.ParkedVehicles = append(result.ParkedVehicles, vehicle)
	}
	return result
}

//GetVehicleByID - This method will retrieve data for a given vehicle, based on the ID sent in the request
func GetVehicleByID(db *gorm.DB, id int) (string, int, float64) {

	var vehicle ParkedVehicle

	err := db.Table("parked_vehicles").Where("is_parked = true and id = ?", id).First(&vehicle).Error
	if err != nil {
		return "", 0, 0
	}
	return vehicle.LicensePlate, vehicle.Duration, vehicle.TotalPaid
}

//PostVehicleEntry - This method is responsible for taking in a 'Vehicle', and storing it's info in the database. Time parked is currently a randomly generated number from 0-24/ 0 will be billed as 1 hour
func PostVehicleEntry(db *gorm.DB, c *viper.Viper, licensePlate string) (ParkedVehicle, error) {

	licensePlate = strings.ToUpper(licensePlate)
	if hasSymbol(licensePlate) {
		errorMsg := fmt.Sprintf("Illegal symbol in License Plate: '%s'. No symbols for you!", licensePlate)
		return ParkedVehicle{}, errors.New(errorMsg)
	}

	var vehicle = ParkedVehicle{LicensePlate: licensePlate}
	randSrc := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(randSrc)
	totalParked := GetVehicles(db)

	if len(totalParked.ParkedVehicles) >= 5 {
		return vehicle, errors.New("Parking lot is full. Please try again later")
	}

	vehicle.IsParked = true
	vehicle.Duration = rand.Intn(24)
	vehicle.TotalPaid = ProcessTotalPaid(c, vehicle.Duration)
	result := db.Create(&vehicle)

	if result.Error != nil {
		return vehicle, result.Error
	}
	return vehicle, result.Error
}

//PostVehiclePayment - This method is responsible for updating the DB's entry for a parked vehicle to be paid & 'removed' from the lot
func PostVehiclePayment(db *gorm.DB, v *viper.Viper, id int) (int64, string, int, float64, error) {

	var result *gorm.DB
	var isRecord bool
	var vehicle ParkedVehicle

	if err := db.Where("id = ?", id).First(&vehicle).Error; err != nil {
		isRecord = false
	} else {
		isRecord = true
	}

	if isRecord {
		result = db.Model(&vehicle).Where("id = ?", id).Update("is_parked", false)
		return result.RowsAffected, vehicle.LicensePlate, vehicle.Duration, vehicle.TotalPaid, nil
	}
	return 0, "N/A", 0, 0, fmt.Errorf("No record found for vehicle with ticket id %d", id)
}

//ProcessTotalPaid - This method is responsible for processing what the Total price is, according to configurable modifiers for price based on hours parked for - refer to config.yml for these rates
func ProcessTotalPaid(c *viper.Viper, duration int) float64 {

	var total float64

	switch {
	case duration <= 1:
		total = c.GetFloat64("config.parking-lot.starting-rate")
		return CustomDecimalRound(total, 0.01)
	case duration <= 3:
		total = (c.GetFloat64("config.parking-lot.starting-rate") * c.GetFloat64("config.parking-lot.three-hour-mod"))
		return CustomDecimalRound(total, 0.01)
	case duration <= 6:
		total = float64((c.GetFloat64("config.parking-lot.starting-rate") * c.GetFloat64("config.parking-lot.six-hour-mod")))
		return CustomDecimalRound(total, 0.01)
	case duration <= 24:
		total = (c.GetFloat64("config.parking-lot.starting-rate") * c.GetFloat64("config.parking-lot.all-day-mod"))
		return CustomDecimalRound(total, 0.01)
	default:
		return 0
	}
}

//CustomDecimalRound - This method is able to round to decimal places using a unit of your choice - e.x. you can send in 0.05 to round to the nearest 5 cents for currency, or 0.01 to round to the nearest cent.
func CustomDecimalRound(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

//hasSymbol - This method is used to check if a license plate submission has any illegal symbols passed in.
func hasSymbol(str string) bool {
	if strings.ContainsAny(str, ",.!\\/|[]{}()-_=+@#$%^&*<>") {
		return true
	}
	return false
}
