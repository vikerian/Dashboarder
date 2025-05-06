package main

// import libs
import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

// types

// Config -> configuration options loaded from environment by method
type Config struct {
	MqttURL         string      `json:"mqtt_url"`
	MqttAUTH        interface{} `json:"mqtt_auth"`
	MqttChannelLIST []string    `json:"mqtt_channels"`
	ValkeyURL       string      `json:"valkey_url"`
	ValkeyAUTH      interface{} `json:"valkey_auth"`
	SiriURL         string      `json:"siridb_url"`
	SiriAUTH        interface{} `json:"siridb_auth"`
	MongoURL        string      `json:"mongo_url"`
	MongoAUTH       interface{} `json:"mongo_auth"`
}

// AuthVAR -> authentication data decoded
type AuthVAR struct {
	User string
	Pass string
}

// NewConfig -> constructor of configuration data instance
func NewConfig() Config {
	cf := new(Config)
	return *cf
}

// LoadConfiguration -> load configuration into our instance from environment (do checks if values exists)
func (cf *Config) LoadConfiguration() (err error) {
	var ok bool
	cf.MqttURL, ok = os.LookupEnv("MQTT_URL")
	if !ok {
		return errors.New("Missing MQTT_URL in environment")
	}

	cf.MqttAUTH, ok = os.LookupEnv("MQTT_AUTH")
	if !ok {
		return errors.New("Missing or nondecodable MQTT_AUTH in environment")
	}

	cf.ValkeyURL, ok = os.LookupEnv("VALKEY_URL")
	if !ok {
		return errors.New("Missing VALKEY_URL in environment")
	}

	cf.ValkeyAUTH, ok = os.LookupEnv("VALKEY_AUTH")
	if !ok {
		return errors.New("Missing VALKEY_AUTH in environment")
	}

	cf.SiriURL, ok = os.LookupEnv("SIRIDB_URL")
	if !ok {
		return errors.New("Missing SIRIDB_URL in environment")
	}

	cf.SiriAUTH, ok = os.LookupEnv("SIRIDB_AUTH")
	if !ok {
		return errors.New("Missing SIRIDB_AUTH in environment")
	}

	cf.MongoURL, ok = os.LookupEnv("MONGODB_URL")
	if !ok {
		return errors.New("Missing MONGODB_URL in environment")
	}

	cf.MongoAUTH, ok = os.LookupEnv("MONGODB_AUTH")
	if !ok {
		return errors.New("Missing MONGODB_AUTHh in environment")
	}
	return
}

// DecodeAuth -> deocde authentication from b64, return decoded data as AuthVAR instance
func DecodeAuth(in string) (av AuthVAR, err error) {
	var buffer bytes.Buffer
	data, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		err = fmt.Errorf("Error on decoding auth: %v", err)
		return
	}
	buffer.Write(data)
	bufstr := buffer.String()
	fmt.Printf("\n %s \n", bufstr)

	return
}
