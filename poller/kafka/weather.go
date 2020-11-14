package kafka

import (
	"encoding/json"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"net/http"
	"poller/model"
)

var apiKey = "a67b20fe5d78b4f268c2fc3d433e9104"

func GetWeatherData(zipCode string) *model.Weather {
	url := "http://api.openweathermap.org/data/2.5/weather?zip=" + zipCode + "&units=imperial&appid=" + apiKey
	resp, err := http.Get(url)
	if err != nil {
		log.Error("Weather api error")
		log.Error(err.Error())
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("Weather api error")
		log.Error(err.Error())
		return nil
	}

	weather := model.Weather{}
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Error("Weather api error")
		log.Error(err.Error())
		return nil
	}

	log.Info("OPEN WEATHER API CALL zipcode: %s", zipCode)

	return &weather
}
