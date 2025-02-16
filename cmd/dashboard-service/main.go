package main

import (
	"dashboarder/internal/api"
	"dashboarder/internal/logger"
	"dashboarder/internal/mongoconfig"
	"fmt"
)

func main() {
	log := logger.GetLogger()
	facade := api.NewAPIFacade()
	mongocfg := mongoconfig.GetMongo()

	// display mongo configuration
	mongocfgstr := fmt.Sprintf("Mongo configuration -> Client: %v, URL: %v", mongocfg.Client, mongocfg.Url)
	log.Info(mongocfgstr)

	// Fetch weather data
	weather, err := facade.FetchData("weather")
	if err != nil {
		log.Error("Failed to fetch weather", "error", err)
		return
	}
	log.Info("Weather data", "weather", weather)

	// Fetch traffic data
	traffic, err := facade.FetchData("traffic")
	if err != nil {
		log.Error("Failed to fetch traffic", "error", err)
		return
	}
	log.Info("Traffic data", "traffic", traffic)

	// Fetch temperature data
	temperature, err := facade.FetchData("temperature")
	if err != nil {
		log.Error("Failed to fetch temperature", "error", err)
		return
	}
	log.Info("Temperature data", "temperature", temperature)

	// Fetch articles
	articles, err := facade.FetchData("article")
	if err != nil {
		log.Error("Failed to fetch articles", "error", err)
		return
	}
	log.Info("Articles", "articles", articles)
}
