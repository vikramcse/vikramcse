package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type weather struct {
	city   string
	apiKey string
}

type weatherInfo struct {
	Request struct {
		Type     string `json:"type"`
		Query    string `json:"query"`
		Language string `json:"language"`
		Unit     string `json:"unit"`
	} `json:"request"`
	Location struct {
		Name           string `json:"name"`
		Country        string `json:"country"`
		Region         string `json:"region"`
		Lat            string `json:"lat"`
		Lon            string `json:"lon"`
		TimezoneID     string `json:"timezone_id"`
		Localtime      string `json:"localtime"`
		LocaltimeEpoch int    `json:"localtime_epoch"`
		UtcOffset      string `json:"utc_offset"`
	} `json:"location"`
	Current struct {
		ObservationTime     string   `json:"observation_time"`
		Temperature         int      `json:"temperature"`
		WeatherCode         int      `json:"weather_code"`
		WeatherIcons        []string `json:"weather_icons"`
		WeatherDescriptions []string `json:"weather_descriptions"`
		WindSpeed           int      `json:"wind_speed"`
		WindDegree          int      `json:"wind_degree"`
		WindDir             string   `json:"wind_dir"`
		Pressure            int      `json:"pressure"`
		Precip              float64  `json:"precip"`
		Humidity            int      `json:"humidity"`
		Cloudcover          int      `json:"cloudcover"`
		Feelslike           int      `json:"feelslike"`
		UvIndex             int      `json:"uv_index"`
		Visibility          int      `json:"visibility"`
		IsDay               string   `json:"is_day"`
	} `json:"current"`
}

func (w *weather) getWeatherInfo() (*weatherInfo, error) {
	info := &weatherInfo{}
	api := "http://api.weatherstack.com/current?access_key=d1ddbfd3540b310d5558148c317b0877" + "&query=" + w.city

	r, err := http.Get(api)
	if err != nil {
		return info, fmt.Errorf("error while fetching information from api %v", err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		return info, fmt.Errorf("error while decoding json %v", err)
	}

	return info, nil
}
