package storage

import (
	"dashboarder/internal/logger"
)

// SiriDBStorage implements Storage for SiriDB
type SiriDBStorage struct{}

func (s *SiriDBStorage) Save(data Data) error {
	logger.GetLogger().Info("Saving to SiriDB", "data", data)
	// Simulate saving to SiriDB
	// In a real implementation, you'd use the SiriDB client to write data
	return nil
}

// MongoStorage implements Storage for MongoDB
type MongoStorage struct{}

func (m *MongoStorage) Save(data Data) error {
	logger.GetLogger().Info("Saving to MongoDB", "data", data)
	// Simulate saving to MongoDB
	// In a real implementation, you'd use the MongoDB client to insert data
	return nil
}

// Data represents the data to be stored
type Data struct {
	Type  string
	Value interface{}
}

// Storage is the interface that all storage must implement
type Storage interface {
	Save(data Data) error
}
