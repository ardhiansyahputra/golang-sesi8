package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

type StatusWeather struct {
	Status Status `json:"status"`
}

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Weather struct {
	Weather int
	Gif     GIF
	Status  Status
}

type GIF struct {
	Id, Name string
	Point    int
}

func newWeather(ticker *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-ticker.C:
			statusWeather := StatusWeather{
				Status: Status{
					Water: rand.Intn(19) + 1,
					Wind:  rand.Intn(19) + 1,
				},
			}

			newStatusWeather, _ := json.Marshal(statusWeather)

			_ = ioutil.WriteFile("weather.json", newStatusWeather, 0666)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func statusWeather(w http.ResponseWriter, r *http.Request) {
	gif := []GIF{{
		Id:    "cnLRDu5HLlQKLOAnhK",
		Name:  "Its Safe !",
		Point: 0,
	}, {
		Id:    "JTy5Vs3SleIhwOKckx",
		Name:  "Warning !!",
		Point: 1,
	}, {
		Id:    "L2I6dMqxS9qaHNGrdp",
		Name:  "Danger !!!",
		Point: 2,
	}}

	jsonFile, _ := os.Open("weather.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var statusWeather StatusWeather

	json.Unmarshal([]byte(byteValue), &statusWeather)

	status := statusWeather.Status
	weather := Weather{
		Weather: 0,
		Status:  statusWeather.Status,
	}

	if (status.Water > 5 && status.Water <= 8) || (status.Wind > 6 && status.Wind <= 15) {
		weather.Weather = 1
	}

	if status.Water > 8 || status.Wind > 15 {
		weather.Weather = 2
	}

	weather.Gif = gif[weather.Weather]

	var t, err = template.ParseFiles("template.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(weather)

	t.Execute(w, weather)
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go newWeather(ticker, quit)

	http.HandleFunc("/", statusWeather)

	fmt.Println("starting web server at http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error running service: ", err)
	}
}
