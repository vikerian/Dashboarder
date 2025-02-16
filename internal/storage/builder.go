package storage

import "errors"

// StorageBuilder decides which storage to use
type StorageBuilder struct{}

func (b *StorageBuilder) Build(data Data) (Storage, error) {
	switch data.Type {
	case "temperature":
		return &SiriDBStorage{}, nil
	case "article":
		return &MongoStorage{}, nil
	default:
		return nil, errors.New("unsupported data type")
	}
}
