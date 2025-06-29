// services/data-router/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"home-dashboard/shared/config"
	"home-dashboard/shared/models"
	mqttlib "home-dashboard/shared/mqtt"

	"github.com/SiriDB/go-siridb-connector"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	// "github.com/transceptor-technology/go-siridb-connector"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataRouter struct {
	mqttClient   *mqttlib.Client
	mongoClient  *mongo.Client
	mongoDB      *mongo.Database
	siridbClient *siridb.Client
	config       *config.Config
}

func NewDataRouter(cfg *config.Config) (*DataRouter, error) {
	// MQTT setup
	mqttConfig := &mqttlib.Config{
		Broker:   cfg.MQTTBroker,
		ClientID: "data-router",
	}

	mqttClient, err := mqttlib.NewClient(mqttConfig)
	if err != nil {
		return nil, fmt.Errorf("mqtt client error: %w", err)
	}

	// MongoDB setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("mongodb connect error: %w", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongodb ping error: %w", err)
	}

	// SiriDB setup
	// NewClient returns a pointer to a new client object.
	// Example hostlist:
	// [][]interface{}{
	//	 {"myhost1", 9000}, 		// hostname/ip and port are required
	//   {"myhost2", 9000, 2},      // an optional integer value can be used as weight
	//								// (default weight is 1)
	//   {"myhost3", 9000, true},   // if true is added as third argument the host
	//								// will be used only when other hosts are not available
	// }
	//
	siridbClient := siridb.NewClient(
		//	cfg.SiriDBHost,
		//	cfg.SiriDBPort,
		cfg.SiriDBUser, // username
		cfg.SiriDBPass, // password
		"home_sensors", // database name
		[][]interface{}{
			{cfg.SiriDBHost, cfg.SiriDBPort},
		},
		nil, // logCh chan string
	)

	siridbClient.Connect()

	return &DataRouter{
		mqttClient:   mqttClient,
		mongoClient:  mongoClient,
		mongoDB:      mongoClient.Database(cfg.MongoDB),
		siridbClient: siridbClient,
		config:       cfg,
	}, nil
}

func (dr *DataRouter) routeMMDecinData(data models.MMDecinData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := dr.mongoDB.Collection("mmdecin_articles")

	// Check if article already exists (by URL)
	filter := map[string]string{"url": data.URL}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("count error: %w", err)
	}

	if count == 0 {
		_, err = collection.InsertOne(ctx, data)
		if err != nil {
			return fmt.Errorf("insert error: %w", err)
		}
		log.Printf("Stored new MMDecin article: %s", data.Title)
	} else {
		log.Printf("Article already exists: %s", data.URL)
	}

	// Publish to processor
	dr.publishToProcessor("mmdecin_new", data)

	return nil
}

func (dr *DataRouter) routeTrafficData(data models.TrafficIncident) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := dr.mongoDB.Collection("traffic_incidents")

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("insert traffic error: %w", err)
	}

	log.Printf("Stored traffic incident: %s at %s", data.Type, data.Location)

	// Publish to processor
	dr.publishToProcessor("traffic_new", data)

	return nil
}

func (dr *DataRouter) routeIoTData(data models.IoTSensorData) error {
	// Store in SiriDB for time-series data
	if dr.siridbClient != nil {
		series := fmt.Sprintf("sensor.%s.%s", data.SensorID, data.SensorType)
		timestamp := data.Timestamp.Unix()

		err := dr.siridbClient.Insert([]siridb.Series{
			{
				Name: series,
				Points: []siridb.Point{
					{Timestamp: timestamp, Value: data.Value},
				},
			},
		})

		if err != nil {
			log.Printf("SiriDB insert error: %v", err)
		} else {
			log.Printf("Stored IoT data: %s = %.2f %s", series, data.Value, data.Unit)
		}
	}

	// Also store metadata in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := dr.mongoDB.Collection("iot_metadata")

	// Upsert sensor metadata
	filter := map[string]string{"sensor_id": data.SensorID}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"sensor_id":    data.SensorID,
			"sensor_type":  data.SensorType,
			"unit":         data.Unit,
			"last_value":   data.Value,
			"last_updated": data.Timestamp,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		log.Printf("MongoDB IoT metadata error: %v", err)
	}

	// Publish to processor
	dr.publishToProcessor("iot_new", data)

	return nil
}

func (dr *DataRouter) publishToProcessor(eventType string, data interface{}) {
	message := models.MQTTMessage{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	if err := dr.mqttClient.Publish("processor/events", message); err != nil {
		log.Printf("Failed to publish to processor: %v", err)
	}
}

func (dr *DataRouter) handleMessage(client mqtt.Client, msg mqtt.Message) {
	var message models.MQTTMessage
	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Received message type: %s", message.Type)

	switch message.Type {
	case "mmdecin_data":
        :x
		if err := mapToStruct(message.Data, &data); err != nil {
			log.Printf("Failed to parse MMDecin data: %v", err)
			return
		}
		if err := dr.routeMMDecinData(data); err != nil {
			log.Printf("Failed to route MMDecin data: %v", err)
		}

	case "traffic_incident":
		var data models.TrafficIncident
		if err := mapToStruct(message.Data, &data); err != nil {
			log.Printf("Failed to parse traffic data: %v", err)
			return
		}
		if err := dr.routeTrafficData(data); err != nil {
			log.Printf("Failed to route traffic data: %v", err)
		}

	case "iot_sensor":
		var data models.IoTSensorData
		if err := mapToStruct(message.Data, &data); err != nil {
			log.Printf("Failed to parse IoT data: %v", err)
			return
		}
		if err := dr.routeIoTData(data); err != nil {
			log.Printf("Failed to route IoT data: %v", err)
		}

	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

func mapToStruct(input interface{}, output interface{}) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, output)
}

func (dr *DataRouter) Start() error {
	// Subscribe to all scraper topics
	topics := []string{
		"scrapers/+",
		"sensors/+",
	}

	for _, topic := range topics {
		if err := dr.mqttClient.Subscribe(topic, dr.handleMessage); err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", topic, err)
		}
		log.Printf("Subscribed to topic: %s", topic)
	}

	log.Println("Data router is running...")

	// Keep the service running
	select {}
}

func (dr *DataRouter) Shutdown() {
	if dr.siridbClient != nil {
		dr.siridbClient.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dr.mongoClient.Disconnect(ctx); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	}

	dr.mqttClient.Disconnect()
}

func main() {
	cfg := config.Load()

	router, err := NewDataRouter(cfg)
	if err != nil {
		log.Fatalf("Failed to create data router: %v", err)
	}
	defer router.Shutdown()

	if err := router.Start(); err != nil {
		log.Fatalf("Failed to start router: %v", err)
	}
}
