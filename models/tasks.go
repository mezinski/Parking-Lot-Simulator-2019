package models

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

//ParkedVehicle ...
type ParkedVehicle struct {
	ID           uint    `json:"id"`
	LicensePlate string  `json:"license_plate"`
	Duration     int     `json:"duration"`
	IsParked     bool    `json:"is_parked"`
	TotalPaid    float64 `json:"total_paid"`
}

//VehicleCollection ...
type VehicleCollection struct {
	ParkedVehicles []ParkedVehicle `json:"parked_vehicles"`
}

//GetVehicleByID ...
func GetVehicleByID(db *gorm.DB, id int) (string, int, float64) {

	var vehicle ParkedVehicle

	err := db.Table("parked_vehicles").Where("is_parked = true and id = ?", id).First(&vehicle).Error
	if err != nil {
		return "", 0, 0
	}
	return vehicle.LicensePlate, vehicle.Duration, vehicle.TotalPaid
}

//GetVehicles ...
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

//PostVehicleEntry ...
func PostVehicleEntry(db *gorm.DB, c *viper.Viper, licensePlate string) (int64, error) {

	var vehicle = ParkedVehicle{LicensePlate: licensePlate}

	randSrc := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(randSrc)

	totalParked := GetVehicles(db)

	if len(totalParked.ParkedVehicles) >= 5 {
		fmt.Println("here")
		return 0, errors.New("Parking lot is full. Please try again later")
	}
	vehicle.IsParked = true
	vehicle.Duration = rand.Intn(24)
	vehicle.TotalPaid = ProcessPayment(c, vehicle.Duration)
	fmt.Println(vehicle.TotalPaid)

	result := db.Create(&vehicle)
	if result.Error != nil {
		return 0, result.Error
	}

	result = db.Save(&vehicle)
	if result.Error != nil {
		return 0, result.Error
	}

	idInt := int64(vehicle.ID)
	return idInt, result.Error
}

//ProcessPayment ...
func ProcessPayment(c *viper.Viper, duration int) float64 {
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

//PostVehicleDuration ...
func PostVehicleDuration(db *gorm.DB, id int, duration int) (int64, string, int, error) {

	var result *gorm.DB
	var isRecord bool
	var vehicle ParkedVehicle

	if err := db.Where("id = ?", id).First(&vehicle).Error; err != nil {
		isRecord = false
	} else {
		isRecord = true
	}

	if isRecord {
		result = db.Model(&vehicle).Where("id = ?", id).Update("duration", duration)
		return result.RowsAffected, vehicle.LicensePlate, vehicle.Duration, nil
	}
	return 0, "N/A", 0, fmt.Errorf("No record found for vehicle with ticket id %d", id)
}

//PostVehiclePayment ...
func PostVehiclePayment(db *gorm.DB, v *viper.Viper, id int) (int64, string, int, float64, error) {

	var result *gorm.DB
	var isRecord bool
	var vehicle ParkedVehicle

	if err := db.Where("id = ?", id).First(&vehicle).Error; err != nil {
		isRecord = false
	} else {
		isRecord = true
	}
	fmt.Println(vehicle)
	if isRecord {
		result = db.Model(&vehicle).Where("id = ?", id).Update("is_parked", false)
		fmt.Println(vehicle)
		return result.RowsAffected, vehicle.LicensePlate, vehicle.Duration, vehicle.TotalPaid, nil
	}
	return 0, "N/A", 0, 0, fmt.Errorf("No record found for vehicle with ticket id %d", id)
}

//DeleteTask ...
//func DeleteTask(db *gorm.DB, id int) (int64, error) {

// var result *gorm.DB
// isRecord := true

// if err := db.Where("id = ?", id).First(&Task{}).Error; err != nil {
// 	isRecord = false
// } else {
// 	isRecord = true
// }

// if isRecord {
// 	result = db.Where("id = ?", id).Delete(&Task{})
// 	return result.RowsAffected, fmt.Errorf("Deleted task with id %d successfully", id)
// }

// return 0, fmt.Errorf("No record found for task with id %d", id)
//}

//CustomDecimalRound ...
func CustomDecimalRound(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
