// services/api-server/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"home-dashboard/shared/config"
	"home-dashboard/shared/models"

	"github.com/gorilla/mux"
	// "github.com/transceptor-technology/go-siridb-connector"
	"github.com/SiriDB/go-siridb-connector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type APIServer struct {
	mongoClient  *mongo.Client
	mongoDB      *mongo.Database
	siridbClient *siridb.Client
	config       *config.Config
	router       *mux.Router
}

func NewAPIServer(cfg *config.Config) (*APIServer, error) {
	// MongoDB setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("mongodb connect error: %w", err)
	}

	// SiriDB setup (optional)
	var siridbClient *siridb.Client
	siridbClient = siridb.NewClient(
		cfg.SiriDBHost,
		cfg.SiriDBPort,
		cfg.SiriDBUser,
		cfg.SiriDBPass,
		"home_sensors",
		false,
		"",
	)

	if err := siridbClient.Connect(); err != nil {
		log.Printf("Warning: SiriDB connection failed: %v", err)
		siridbClient = nil
	}

	api := &APIServer{
		mongoClient:  mongoClient,
		mongoDB:      mongoClient.Database(cfg.MongoDB),
		siridbClient: siridbClient,
		config:       cfg,
		router:       mux.NewRouter(),
	}

	api.setupRoutes()
	return api, nil
}

func (a *APIServer) setupRoutes() {
	// Middleware
	a.router.Use(corsMiddleware)
	a.router.Use(jsonMiddleware)

	// Routes
	a.router.HandleFunc("/api/health", a.healthCheck).Methods("GET")

	// MMDecin routes
	a.router.HandleFunc("/api/mmdecin", a.getMMDecinArticles).Methods("GET")
	a.router.HandleFunc("/api/mmdecin/{id}", a.getMMDecinArticle).Methods("GET")

	// Traffic routes
	a.router.HandleFunc("/api/traffic", a.getTrafficIncidents).Methods("GET")
	a.router.HandleFunc("/api/traffic/active", a.getActiveTrafficIncidents).Methods("GET")

	// IoT routes
	a.router.HandleFunc("/api/sensors", a.getSensors).Methods("GET")
	a.router.HandleFunc("/api/sensors/{id}/data", a.getSensorData).Methods("GET")

	// Statistics
	a.router.HandleFunc("/api/stats", a.getStatistics).Methods("GET")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (a *APIServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now(),
		"services": map[string]bool{
			"mongodb": a.mongoClient.Ping(context.Background(), nil) == nil,
			"siridb":  a.siridbClient != nil,
		},
	}
	json.NewEncoder(w).Encode(health)
}

func (a *APIServer) getMMDecinArticles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse query parameters
	limit := 20
	skip := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	if s := r.URL.Query().Get("skip"); s != "" {
		if val, err := strconv.Atoi(s); err == nil && val >= 0 {
			skip = val
		}
	}

	collection := a.mongoDB.Collection("mmdecin_articles")

	// Sort by timestamp descending (newest first)
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(skip))

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var articles []models.MMDecinData
	if err := cursor.All(ctx, &articles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(articles)
}

func (a *APIServer) getMMDecinArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := a.mongoDB.Collection("mmdecin_articles")

	var article models.MMDecinData
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (a *APIServer) getTrafficIncidents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse time range
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	filter := bson.M{}

	if fromStr != "" || toStr != "" {
		timeFilter := bson.M{}

		if fromStr != "" {
			if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
				timeFilter["$gte"] = from
			}
		}

		if toStr != "" {
			if to, err := time.Parse(time.RFC3339, toStr); err == nil {
				timeFilter["$lte"] = to
			}
		}

		if len(timeFilter) > 0 {
			filter["timestamp"] = timeFilter
		}
	}

	collection := a.mongoDB.Collection("traffic_incidents")

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var incidents []models.TrafficIncident
	if err := cursor.All(ctx, &incidents); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(incidents)
}

func (a *APIServer) getActiveTrafficIncidents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Active incidents = last 24 hours
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": time.Now().Add(-24 * time.Hour),
		},
	}

	collection := a.mongoDB.Collection("traffic_incidents")

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var incidents []models.TrafficIncident
	if err := cursor.All(ctx, &incidents); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(incidents)
}

func (a *APIServer) getSensors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := a.mongoDB.Collection("iot_metadata")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var sensors []map[string]interface{}
	if err := cursor.All(ctx, &sensors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sensors)
}

func (a *APIServer) getSensorData(w http.ResponseWriter, r *http.Request) {
	if a.siridbClient == nil {
		http.Error(w, "SiriDB not available", http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	sensorID := vars["id"]

	// Parse query parameters
	sensorType := r.URL.Query().Get("type")
	if sensorType == "" {
		sensorType = "temperature" // default
	}

	// Time range (default: last 24 hours)
	hours := 24
	if h := r.URL.Query().Get("hours"); h != "" {
		if val, err := strconv.Atoi(h); err == nil && val > 0 && val <= 168 {
			hours = val
		}
	}

	series := fmt.Sprintf("sensor.%s.%s", sensorID, sensorType)
	from := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
	to := time.Now().Unix()

	query := fmt.Sprintf("select * from '%s' between %d and %d", series, from, to)

	result, err := a.siridbClient.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Transform SiriDB result to JSON-friendly format
	data := map[string]interface{}{
		"sensor_id": sensorID,
		"type":      sensorType,
		"from":      from,
		"to":        to,
		"points":    []map[string]interface{}{},
	}

	// Parse result (this depends on SiriDB response format)
	if series, ok := result[series]; ok {
		if points, ok := series.([][]interface{}); ok {
			for _, point := range points {
				if len(point) >= 2 {
					data["points"] = append(data["points"].([]map[string]interface{}), map[string]interface{}{
						"timestamp": point[0],
						"value":     point[1],
					})
				}
			}
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (a *APIServer) getStatistics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats := map[string]interface{}{}

	// MMDecin article count
	articleCount, _ := a.mongoDB.Collection("mmdecin_articles").CountDocuments(ctx, bson.M{})
	stats["mmdecin_articles"] = articleCount

	// Traffic incidents (last 24h and total)
	trafficCollection := a.mongoDB.Collection("traffic_incidents")
	totalIncidents, _ := trafficCollection.CountDocuments(ctx, bson.M{})
	activeIncidents, _ := trafficCollection.CountDocuments(ctx, bson.M{
		"timestamp": bson.M{"$gte": time.Now().Add(-24 * time.Hour)},
	})

	stats["traffic"] = map[string]interface{}{
		"total":  totalIncidents,
		"active": activeIncidents,
	}

	// Sensor count
	sensorCount, _ := a.mongoDB.Collection("iot_metadata").CountDocuments(ctx, bson.M{})
	stats["sensors"] = sensorCount

	json.NewEncoder(w).Encode(stats)
}

func (a *APIServer) Start() error {
	addr := ":" + a.config.APIPort
	log.Printf("Starting API server on %s", addr)

	return http.ListenAndServe(addr, a.router)
}

func (a *APIServer) Shutdown() {
	if a.siridbClient != nil {
		a.siridbClient.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.mongoClient.Disconnect(ctx); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	}
}

func main() {
	cfg := config.Load()

	server, err := NewAPIServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create API server: %v", err)
	}
	defer server.Shutdown()

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
