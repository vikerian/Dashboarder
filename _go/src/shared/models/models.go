package models

import "time"

// MMDecinData represents scraped data from mmdecin.cz
type MMDecinData struct {
	ID        string    `json:"id" bdon:"_id,omitempty"`
	Title     string    `json:"title" "bdon:"title"`
	Content   string    `json:"content" bson:"content"`
	URL       string    `json:"url" bson:"url"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// TrafficIncident represents traffic data for Ustecky kraj
type TrafficIncident struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Type        string    `json:"type" bson:"type"`
	Location    string    `json:"location" bson:"location"`
	Description string    `json:"description" bson:"description"`
	Severity    string    `json:"severity" bson:"severity"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
	Lat         float64   `json:"lat,omitempty" bson:"lat,omitempty"`
	Lng         float64   `json:"lng,omitempty" bson:"lng,omitempty"`
}

// IotSensorData for SiriDB
type IoTSensorData struct {
	SensorID   string    `json:"sensor_id"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
	SensorType string    `json:"sensor_type"` // temperature, humidity, etc
}

// MQTTMessage wrapper for all messages
type MQTTMessage struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}
