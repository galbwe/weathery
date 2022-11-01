package main

import (
	"fmt"
	"log"

	"example.com/weathery"
)

func main() {

	log.SetPrefix("weathery: ")
	log.SetFlags(0)

	sanFrancisco, _ := weathery.GetAccuWeatherData("san francisco")
	denver, _ := weathery.GetAccuWeatherData("denver")
	kansasCity, _ := weathery.GetAccuWeatherData("kansas city")
	newYorkCity, _ := weathery.GetAccuWeatherData("new york")
	fmt.Println(sanFrancisco)
	fmt.Println(denver)
	fmt.Println(kansasCity)
	fmt.Println(newYorkCity)
}
