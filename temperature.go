package main

import (
	"errors"
	"math"
	"strconv"

	"github.com/martinlindhe/unit"
)

func getTemperatureValue(message string) []string {
	submatchall := numberRegex.FindAllString(message, -1)
	return submatchall
}

func getFarenheit(message string) ([]string, []float64, error) {

	var calculations []float64
	temperatures := getTemperatureValue(message)
	if len(temperatures) == 0 {
		err := errors.New("no temperatures found in the message")
		return temperatures, calculations, err
	}
	for _, temperature := range temperatures {
		temp, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			return temperatures, calculations, err
		}
		c := unit.FromCelsius(temp)
		calculations = append(calculations, math.Round(c.Fahrenheit()*100)/100)

	}
	return temperatures, calculations, nil
}

func getCelcius(message string) ([]string, []float64, error) {
	var calculations []float64
	temperatures := getTemperatureValue(message)
	if len(temperatures) == 0 {
		err := errors.New("no temperatures found in the message")
		return temperatures, calculations, err
	}
	for _, temperature := range temperatures {
		temp, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			return temperatures, calculations, err
		}
		f := unit.FromFahrenheit(temp)
		calculations = append(calculations, math.Round(f.Celsius()*100)/100)

	}
	return temperatures, calculations, nil
}
