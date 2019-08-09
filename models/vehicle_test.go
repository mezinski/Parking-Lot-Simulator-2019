package models

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

/*
Foreword

Just a note - I am aware that these tests only encapsulate a small amount
of methods that are present in the application, and this is for reason.

I would like to use this just as a show of knowledge for testing, and that I'm
versed in how tests work. I did see that tests were an optional choice for the
technical test, and would much prefer any discussion or time put into this to
go to me asking you about how you do your tests, and what your TDD looks like.
*/

func TestVehicles(t *testing.T) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("../")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			t.Error(fmt.Errorf("Config file was not found: %s", err.Error()))
		} else {
			t.Error(fmt.Errorf("Config file was found but other error is present: %s", err.Error()))
		}
	}

	//-------------------------------------------------------------\\
	// This is for potential use in testing REST API components
	//
	// db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", v.GetString("config.database.host"), v.GetString("config.database.port"), v.GetString("config.database.user"), v.GetString("config.database.dbname")))
	// if err != nil {
	// 	t.Error(err)
	// }

	//Testing ProcessTotalPrice method
	//Passing in all Durations for parking time-slots, and checking to see if the price is accurate when calculated
	t.Run("ProcessTotalPrice - Less than 1 hour", testProcessTotalPrice(v, 0, 3))
	t.Run("ProcessTotalPrice - 1 Hour Parking", testProcessTotalPrice(v, 1, 3))
	t.Run("ProcessTotalPrice - 2-3 Hour Parking", testProcessTotalPrice(v, 3, 4.5))
	t.Run("ProcessTotalPrice - 4-6 Hour Parking", testProcessTotalPrice(v, 6, 6.75))
	t.Run("ProcessTotalPrice - 24 Hour Parking", testProcessTotalPrice(v, 24, 10.13))

	//Testing CustomDecimalRound method
	//Passing in all prices for parking time-slots
	t.Run("CustomDecimalRound - Less than 1 or equal to 1 Hour Price", testCustomDecimalRound(3, 0.01, 3))
	t.Run("CustomDecimalRound - More than 1 or equal to 3 Hour Price", testCustomDecimalRound(4.5, 0.01, 4.5))
	t.Run("CustomDecimalRound - More than 3 or equal to 6 Hour Price", testCustomDecimalRound(6.75, 0.01, 6.75))
	t.Run("CustomDecimalRound - More than 6 or equal to 24 Hour Price", testCustomDecimalRound(10.125, 0.01, 10.13))

	//Testing hasIllegalSymbol method
	//Passing in all of the symbols not allowed to make sure the method returns true
	t.Run("hasIllegalSymbol - ',' character", testHasIllegalSymbol(",", true))
	t.Run("hasIllegalSymbol - '.' character", testHasIllegalSymbol(".", true))
	t.Run("hasIllegalSymbol - '!' character", testHasIllegalSymbol("!", true))
	t.Run("hasIllegalSymbol - '\\' character", testHasIllegalSymbol("\\", true))
	t.Run("hasIllegalSymbol - '/' character", testHasIllegalSymbol("/", true))
	t.Run("hasIllegalSymbol - '|' character", testHasIllegalSymbol("|", true))
	t.Run("hasIllegalSymbol - '[' character", testHasIllegalSymbol("[", true))
	t.Run("hasIllegalSymbol - ']' character", testHasIllegalSymbol("]", true))
	t.Run("hasIllegalSymbol - '{' character", testHasIllegalSymbol("{", true))
	t.Run("hasIllegalSymbol - '}' character", testHasIllegalSymbol("}", true))
	t.Run("hasIllegalSymbol - '(' character", testHasIllegalSymbol("(", true))
	t.Run("hasIllegalSymbol - ')' character", testHasIllegalSymbol(")", true))
	t.Run("hasIllegalSymbol - '-' character", testHasIllegalSymbol("-", true))
	t.Run("hasIllegalSymbol - '_' character", testHasIllegalSymbol("_", true))
	t.Run("hasIllegalSymbol - '=' character", testHasIllegalSymbol("=", true))
	t.Run("hasIllegalSymbol - '+' character", testHasIllegalSymbol("+", true))
	t.Run("hasIllegalSymbol - '@' character", testHasIllegalSymbol("@", true))
	t.Run("hasIllegalSymbol - '#' character", testHasIllegalSymbol("#", true))
	t.Run("hasIllegalSymbol - '$' character", testHasIllegalSymbol("$", true))
	t.Run("hasIllegalSymbol - '%' character", testHasIllegalSymbol("%", true))
	t.Run("hasIllegalSymbol - '^' character", testHasIllegalSymbol("^", true))
	t.Run("hasIllegalSymbol - '&' character", testHasIllegalSymbol("&", true))
	t.Run("hasIllegalSymbol - '*' character", testHasIllegalSymbol("*", true))
	t.Run("hasIllegalSymbol - '<' character", testHasIllegalSymbol("<", true))
	t.Run("hasIllegalSymbol - '>' character", testHasIllegalSymbol(">", true))
}

func testProcessTotalPrice(v *viper.Viper, duration int, expected float64) func(t *testing.T) {
	return func(t *testing.T) {
		actual := ProcessTotalPrice(v, duration)
		if actual != expected {
			t.Error(fmt.Sprintf("Test failed: Expected %f but got %f.", expected, actual))
		}
	}
}

func testCustomDecimalRound(total float64, unit float64, expected float64) func(t *testing.T) {
	return func(t *testing.T) {
		roundedTotal := CustomDecimalRound(total, unit)
		if roundedTotal != expected {
			t.Error(fmt.Sprintf("Test Failed - Expected %f but got %f", expected, roundedTotal))
		}
	}
}

func testHasIllegalSymbol(str string, expected bool) func(t *testing.T) {
	return func(t *testing.T) {
		actual := hasIllegalSymbol(str)
		if actual != expected {
			t.Error(fmt.Sprintf("Test Failed - Expected %v but got %v", expected, actual))
		}
	}
}
