package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

//DbConnection - This method will initialize a GORM DB object with the database information.
func DbConnection(v *viper.Viper) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", v.GetString("config.database.host"), v.GetString("config.database.port"), v.GetString("config.database.user"), v.GetString("config.database.dbname")))
	if err != nil {
		panic(err)
	}
	return db, err
}
