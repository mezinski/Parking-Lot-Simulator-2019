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

//Vehicle - This is the struct for a 'Parked Vehicle'. This is all the information we store and utilize for a vehicle parking with this application.
type Vehicle struct {
	ID           uint    `json:"id"`
	LicensePlate string  `json:"license_plate"`
	Duration     int     `json:"duration"`
	IsParked     bool    `json:"is_parked"`
	TotalPaid    float64 `json:"total_paid"`
}

//TableName - Sets the table name for Vehicle objects
func (Vehicle) TableName() string {
	return "parked_vehicles"
}

//GetVehicles - This method will return all of the vehicles currently parked in the lot
func GetVehicles(db *gorm.DB) []Vehicle {
	vehicles := make([]Vehicle, 0)
	db.Where("is_parked = true").Find(&vehicles)
	return vehicles
}

//GetVehicleByID - This method will retrieve data for a given vehicle, based on the ID sent in the request
func GetVehicleByID(db *gorm.DB, id int) Vehicle {
	var vehicle Vehicle

	err := db.Table("parked_vehicles").Where("is_parked = true and id = ?", id).First(&vehicle).Error
	if err != nil {
		return Vehicle{}
	}
	return vehicle
}

//CreateVehicle - This method is responsible for taking in a 'Vehicle', and storing it's info in the database. Time parked is currently a randomly generated number from 0-24/ 0 will be billed as 1 hour
func CreateVehicle(db *gorm.DB, c *viper.Viper, licensePlate string) (Vehicle, error) {
	totalParked := GetVehicles(db)
	if len(totalParked) >= c.GetInt("config.parking-lot.max-occupancy") {
		return Vehicle{}, errors.New("Parking lot is full. Please try again later")
	}

	licensePlate = strings.ToUpper(licensePlate)
	if hasIllegalSymbol(licensePlate) {
		errorMsg := fmt.Sprintf("Illegal symbol in License Plate: '%s'. No symbols for you!", licensePlate)
		return Vehicle{}, errors.New(errorMsg)
	}

	var vehicle = Vehicle{LicensePlate: licensePlate}
	randSrc := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(randSrc)

	vehicle.IsParked = true
	vehicle.Duration = rand.Intn(24)
	vehicle.TotalPaid = ProcessTotalPrice(c, vehicle.Duration)

	result := db.Create(&vehicle)
	if result.Error != nil {
		return vehicle, result.Error
	}
	return vehicle, result.Error
}

//PostVehiclePayment - This method is responsible for updating the DB's entry for a parked vehicle to be paid & 'removed' from the lot
func (v *Vehicle) PostVehiclePayment(db *gorm.DB, id int) error {
	var vehicle Vehicle
	if err := db.Where("id = ?", id).First(&v).Error; err != nil {
		return fmt.Errorf("No record found for vehicle with ticket id %d", id)
	}
	db.Model(&vehicle).Where("id = ?", id).Update("is_parked", false)
	return nil
}

//ProcessTotalPrice - This method is responsible for processing what the Total price is, according to configurable modifiers for price based on hours parked for - refer to config.yml for these rates
func ProcessTotalPrice(c *viper.Viper, duration int) float64 {
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
func hasIllegalSymbol(str string) bool {
	if strings.ContainsAny(str, ",.!\\/|[]{}()-_=+@#$%^&*<>") {
		return true
	}
	return false
}
