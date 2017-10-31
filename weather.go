package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type weatherProvider interface {
	temperature(country string, city string) (float64, error)
}

type openWeatherMap struct {
	apiKey string
}

func (w openWeatherMap) temperature(country string, city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.apiKey + "&q=" + city)
	if err != nil {
		log.Printf("Error %s", err)
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Name string `json:"name"`
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Printf("Error %s", err)
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}

type weatherUnderground struct {
	apiKey string
}

func (w weatherUnderground) temperature(country string, city string) (float64, error) {
	resp, err := http.Get("http://api.wunderground.com/api/" + w.apiKey + "/conditions/q/" + country + "/" + city + ".json")
	if err != nil {
		log.Printf("Error %s", err)
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Observation struct {
			Celsius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Printf("Error %s", err)
		return 0, err
	}

	kelvin := d.Observation.Celsius + 273.15
	log.Printf("weatherUnderground:: %s: %.2f", city, kelvin)
	return kelvin, nil
}

type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	sum := 0.0
	for _, provider := range w {
		k, err := provider.temperature(city)
		if err != nil {
			return 0, err
		}

		sum += k
	}

	return sum / float64(len(w)), nil
}
