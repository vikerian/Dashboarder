package main

import (
	"dashboarder/internal/api"
	"dashboarder/internal/logger"
)

func main() {
	log := logger.GetLogger()
	facade := api.NewAPIFacade()

	weather, err := facade.FetchData("weather")
	if err != nil {
		log.Error("Failed to fetch weather", "error", err)
		return
	}
	log.Info("Weather data", "weather", weather)
}
