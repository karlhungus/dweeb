package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("need tokens openweather, wunderground")
		os.Exit(1)
	}
	token1 := os.Args[1]
	token2 := os.Args[2]

	mw := multiWeatherProvider{
		openWeatherMap{apiKey: token1},
		weatherUnderground{apiKey: token2},
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		slice := strings.SplitN(r.URL.Path, "/", 4)
		city := slice[3]
		country := slice[2]
		log.Printf("country: %s, city: %s", country, city)

		temp, err := mw.temperature(country, city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset-utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"city":   city,
			"temp":   temp,
			"temp_c": temp - 273.15,
			"took":   time.Since(begin).String(),
		})
	})
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}
