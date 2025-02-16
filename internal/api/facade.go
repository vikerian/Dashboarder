package api

import (
	"dashboarder/internal/logger"
	"errors"
)

// APIFacade provides a simplified interface to fetch data
type APIFacade struct{}

func NewAPIFacade() *APIFacade {
	return &APIFacade{}
}

// FetchData fetches data based on the type
func (f *APIFacade) FetchData(dataType string) (interface{}, error) {
	switch dataType {
	case "weather":
		return fetchWeather(), nil
	case "traffic":
		return fetchTraffic(), nil
	case "temperature":
		return fetchTemperature(), nil
	case "article":
		return fetchArticles(), nil
	default:
		return nil, errors.New("unsupported data type")
	}
}

// Example functions for fetching data
func fetchWeather() interface{} {
	logger.GetLogger().Info("Fetching weather data")
	// Mock data for weather
	return map[string]interface{}{
		"temperature": 22.5,
		"condition":   "Sunny",
	}
}

func fetchTraffic() interface{} {
	logger.GetLogger().Info("Fetching traffic data")
	// Mock data for traffic
	return map[string]interface{}{
		"bus":       "201",
		"station":   "Kamenicka",
		"direction": "Chrochvice",
		"delay":     "5 minutes",
	}
}

func fetchTemperature() interface{} {
	logger.GetLogger().Info("Fetching temperature data")
	// Mock data for temperature
	return 22.5
}

func fetchArticles() interface{} {
	logger.GetLogger().Info("Fetching articles")
	// Mock data for articles
	return []string{"Article 1", "Article 2"}
}
