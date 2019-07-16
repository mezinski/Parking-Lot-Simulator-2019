package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

//Task struct
type Task struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

//ParkedVehicle ...
type ParkedVehicle struct {
	ID           uint   `json:"id"`
	LicensePlate string `json:"license_plate"`
	Duration     uint   `json:"duration"`
	IsParked     bool   `json:"is_parked"`
	TotalPaid    uint   `json:"total_paid"`
}

//TaskCollection ...
type TaskCollection struct {
	Tasks []Task `json:"items"`
}

//VehicleCollection ...
type VehicleCollection struct {
	ParkedVehicles []ParkedVehicle `json:"parked_vehicles"`
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

//GetTasks ...
func GetTasks(db *gorm.DB) TaskCollection {
	rows, err := db.Table("tasks").Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	result := TaskCollection{}
	for rows.Next() {
		task := Task{}
		err2 := rows.Scan(&task.ID, &task.Name)
		if err2 != nil {
			panic(err2)
		}
		result.Tasks = append(result.Tasks, task)
	}
	return result
}

//PostVehicleEntry ...
func PostVehicleEntry(db *gorm.DB, licensePlate string) (int64, error) {

	var vehicle = ParkedVehicle{LicensePlate: licensePlate}

	totalParked := GetVehicles(db)

	if len(totalParked.ParkedVehicles) >= 5 {
		fmt.Println("here")
		return 0, fmt.Errorf("lot is full")
	}
	vehicle.IsParked = true

	result := db.Create(&vehicle)
	if result.Error != nil {
		panic(result.Error)
	}

	result = db.Save(&vehicle)
	if result.Error != nil {
		panic(result.Error)
	}

	idInt := int64(vehicle.ID)
	return idInt, result.Error
}

//UpdateParkingDuration ...
func UpdateParkingDuration(db *gorm.DB) (int64, error) {
	return 0, nil
}

//PutTask ...
func PutTask(db *gorm.DB, name string) (int64, error) {

	var task = Task{Name: name}

	result := db.Create(&task)
	if result.Error != nil {
		panic(result.Error)
	}

	result = db.Save(&task)
	if result.Error != nil {
		panic(result.Error)
	}

	idInt := int64(task.ID)
	return idInt, result.Error
}

//PostVehiclePayment ...
func PostVehiclePayment(db *gorm.DB, id int) (int64, string, uint, uint, error) {

	var result *gorm.DB
	var isRecord bool
	var vehicle ParkedVehicle

	if err := db.Where("id = ?", id).First(&vehicle).Error; err != nil {
		isRecord = false
	} else {
		isRecord = true
	}

	if isRecord {
		result = db.Where("id = ?", id).Update("duration = ?, total_paid = ?, is_parked = false", 1, 3)
		return result.RowsAffected, vehicle.LicensePlate, vehicle.Duration, vehicle.TotalPaid, fmt.Errorf("Thank you for choosing sketchypark. Payment has been processed by your VISA, Card no. %d", viper.GetInt("config.user-data.credit-card-no"))
	}

	return 0, "N/A", 0, 0, fmt.Errorf("No record found for vehicle with ticket id %d", id)
}

//DeleteTask ...
func DeleteTask(db *gorm.DB, id int) (int64, error) {

	var result *gorm.DB
	isRecord := true

	if err := db.Where("id = ?", id).First(&Task{}).Error; err != nil {
		isRecord = false
	} else {
		isRecord = true
	}

	if isRecord {
		result = db.Where("id = ?", id).Delete(&Task{})
		return result.RowsAffected, fmt.Errorf("Deleted task with id %d successfully", id)
	}

	return 0, fmt.Errorf("No record found for task with id %d", id)
}
