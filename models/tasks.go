package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//Task struct
type Task struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

//TaskCollection ...
type TaskCollection struct {
	Tasks []Task `json:"items"`
}

func main() {

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
