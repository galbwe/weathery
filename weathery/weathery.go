package weathery

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type AccuWeatherData struct {
	temperature               int
	temperatureUnit           string
	airQualityMeasurement     int
	airQualityMeasurementUnit string
	airQualityDescription     string
	date                      string
	time                      string
}

func makeAccuWeatherRequest() (*http.Response, error) {
	url := "https://www.accuweather.com/en/us/denver/80203/weather-forecast/347810"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authority", "www.accuweather.com")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("accept-language", "en-US,en;q=0.5")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("sec-gpc", "1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36")

	return client.Do(req)
}

func GetAccuWeatherHtml() (string, error) {
	resp, err := makeAccuWeatherRequest()

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(body)
	unencoded, err := io.ReadAll(reader)
	return string(unencoded[:]), err
}

func GetAccuWeatherData() (AccuWeatherData, error) {
	data := new(AccuWeatherData)
	resp, _ := makeAccuWeatherRequest()
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	temp := doc.Find(".cur-con-weather-card .temp").First().Text()

	// find temperature and temperature unit in temp string
	tempRe := regexp.MustCompile(`\d{1,3}`)
	unitRe := regexp.MustCompile(`[FC]`)
	ts := tempRe.FindString(temp)
	if ts != "" {
		ti, _ := strconv.Atoi(ts)
		data.temperature = ti
	}
	u := unitRe.FindString(temp)
	if u != "" {
		data.temperatureUnit = u
	}

	// find time of weather measurement
	timeRaw := doc.Find(".cur-con-weather-card .cur-con-weather-card__subtitle").Text()
	timeRe := regexp.MustCompile(`\d{1,2}:\d{1,2}\s(A|P)M`)
	timeParsed := timeRe.FindString(timeRaw)
	if timeParsed != "" {
		data.time = timeParsed
	}

	// find air-quality number and unit
	airQualityContent := doc.Find(".air-quality-content")

	// parse aq-number
	airQualityNumber := airQualityContent.Find(".aq-number").Text()
	aqnRe := regexp.MustCompile(`\d{1,3}`)
	aqn := aqnRe.FindString(airQualityNumber)
	if aqn != "" {
		data.airQualityMeasurement, _ = strconv.Atoi(aqn)
	}
	// parse aq-unit
	airQualityUnit := airQualityContent.Find(".aq-unit").Text()
	aquRe := regexp.MustCompile(`\w{1,4}`)
	aqu := aquRe.FindString(airQualityUnit)
	if aqu != "" {
		data.airQualityMeasurementUnit = aqu
	}
	// parse aq-description
	aqDescRaw := airQualityContent.Find(".category-text").Text()
	aqDescRe := regexp.MustCompile(`[a-zA-Z]+`)
	aqDescParsed := aqDescRe.FindString(aqDescRaw)
	if aqDescParsed != "" {
		data.airQualityDescription = aqDescParsed
	}
	// parse date
	dateRaw := airQualityContent.Find(".date").Text()
	dateRe := regexp.MustCompile(`\d{1,2}/\d{1,2}`)
	dateParsed := dateRe.FindString(dateRaw)
	if dateParsed != "" {
		data.date = dateParsed
	}
	return *data, nil
}
