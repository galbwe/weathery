package main

import (
	"fmt"
	"log"

	"example.com/weathery"
)

func main() {

	log.SetPrefix("weathery: ")
	log.SetFlags(0)

	data, _ := weathery.GetAccuWeatherData()
	fmt.Println(data)
}
