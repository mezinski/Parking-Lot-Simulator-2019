package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

//DbConnection - This method will initialize a GORM DB object with the database information.
func DbConnection(c *viper.Viper) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", c.GetString("config.database.host"), c.GetString("config.database.port"), c.GetString("config.database.user"), c.GetString("config.database.dbname")))
	if err != nil {
		panic(err)
	}
	return db, err
}
