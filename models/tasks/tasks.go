package tasks

//Task struct
type Task struct {
	id   uint `gorm:"primary_key"`
	name string
}