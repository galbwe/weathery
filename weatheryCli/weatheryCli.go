package main

import (
	"fmt"
	"log"

	"example.com/weathery"
)

func main() {

	log.SetPrefix("weathery: ")
	log.SetFlags(0)

	temp, _ := weathery.GetAccuWeatherTemperature()
	fmt.Println(temp)
}
