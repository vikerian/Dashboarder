package main

import (
	//"google.golang.org/genproto/googleapis/type/datetime"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const OpenWeatherApiURL = "https://api.openweathermap.org/data/2.5/weather"

type OpenWeatherAPI struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

type OpenWeatherData struct {
	Coord struct {
		Lon float64 `json:"lon,omitempty"`
		Lat float64 `json:"lat,omitempty"`
	} `json:"coord,omitempty"`
	Weather []struct {
		ID          int         `json:"id,omitempty"`
		Main        interface{} `json:"main,omitempty"`
		Description string      `json:"description,omitempty"`
		Icon        string      `json:"icon,omitempty"`
	} `json:"weather,omitempty"`
	Base interface{} `json:"base,omitempty"`
	Main struct {
		Temp      float64 `json:"temp,omitempty"`
		FeelsLike float64 `json:"feels_like,omitempty"`
		TempMin   float64 `json:"temp_min,omitempty"`
		TempMax   float64 `json:"temp_max,omitempty"`
		Pressure  float64 `json:"pressure,omitempty"`
		Humidity  float64 `json:"humidity,omitempty"`
	} `json:"main,omitempty"`
	Visibility int64 `json:"visibility,omitempty"`
	Wind       struct {
		Speed float64 `json:"speed,omitempty"`
		Deg   float64 `json:"deg,omitempty"`
		Gust  float64 `json:"gust,omitempty"`
	} `json:"wind,omitempty"`
	Clouds struct {
		All float64 `json:"all,omitempty"`
	} `json:"clouds,omitempty"`
	Rain struct {
		OneHour   float64 `json:"1h,omitempty"`
		ThreeHour float64 `json:"3h,omitempty"`
	} `json:"rain,omitempty"`
	Snow struct {
		OneHour   float64 `json:"1h,omitempty"`
		ThreeHour float64 `json:"3h,omitempty"`
	} `json:"snow,omitempty"`
	DateTime int64    `json:"dt,omitempty"`
	Sys      struct { //internal for openweathermap
		Type    int64  `json:"type,omitempty"`    //internal for openweathermap
		ID      int    `json:"id,omitempty"`      //internal for openweathermap
		Message string `json:"message,omitempty"` //internal for openweathermap
		Country string `json:"country,omitempty"` //internal for openweathermap
		Sunrise int64  `json:"sunrise,omitempty"` //internal for openweathermap
		Sunset  int64  `json:"sunset,omitempty"`
	} `json:"sys,omitempty"`
	Timezone int64       `json:"timezone,omitempty"` // shift in seconds
	ID       int         `json:"id,omitempty"`       // city id
	Name     string      `json:"name,omitempty"`     // city name
	Cod      interface{} `json:"cod,omitempty"`      //internal for openweathermap
}

func NewOpenWeatherAPI(token string) *OpenWeatherAPI {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &OpenWeatherAPI{
		BaseURL: OpenWeatherApiURL,
		Token:   token,
		Client: &http.Client{
			Timeout:   time.Second * 5,
			Transport: tr,
		},
	}
}

func (ow *OpenWeatherAPI) ReadWeatherData(longitude float64, latitude float64, id int, name string) (*OpenWeatherData, error) {
	// construct query url
	url := fmt.Sprintf("%s?units=metric&lang=cz&lat=%.2f&lon=%.2f&appid=%s", ow.BaseURL, latitude, longitude, ow.Token)

	// create session, setup headers, etc...
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "DCUK-OpenWeatherClient/0.01")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Pragma", "no-cache")

	// do request
	resp, err := ow.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response
	//var respData []byte
	if resp.StatusCode != http.StatusOK {
		//              respData, err = io.ReadAll(resp.Body)
		//              if err != nil {
		//                      return nil, err
		//              }
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	// decode data
	data := new(OpenWeatherData)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}
	data.ID = id
	data.Name = name
	return data, nil
}
