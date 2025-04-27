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
type Config struct {
	MqttUrl         string      `json:"mqtt_url"`
	MqttAuth        interface{} `json:"mqtt_auth"`
	MqttChannelLIST []string    `json:"mqtt_channels"`
	ValkeyUrl       string      `json:"valkey_url"`
	ValkeyAuth      interface{} `json:"valkey_auth"`
	SiridbUrl       string      `json:"siridb_url"`
	SiridbAuth      interface{} `json:"siridb_auth"`
}

type AuthVar struct {
	User string
	Pass string
}

// constructors
func NewConfig() Config {
	cf := new(Config)
	return *cf
}

// methods
func (cf *Config) LoadConfiguration() (err error) {
	var ok bool
	cf.MqttUrl, ok = os.LookupEnv("MQTT_URL")
	if !ok {
		return errors.New("Missing MQTT_URL in environment")
	}

	cf.MqttAuth, ok = os.LookupEnv("MQTT_AUTH")
	if !ok {
		return errors.New("Missing or nondecodable MQTT_AUTH in environment")
	}

	cf.ValkeyUrl, ok = os.LookupEnv("VALKEY_URL")
	if !ok {
		return errors.New("Missing VALKEY_URL in environment")
	}

	cf.ValkeyAuth, ok = os.LookupEnv("VALKEY_AUTH")
	if !ok {
		return errors.New("Missing VALKEY_AUTH in environment")
	}

	cf.SiridbUrl, ok = os.LookupEnv("SIRIDB_URL")
	if !ok {
		return errors.New("Missing SIRIDB_URL in environment")
	}

	cf.SiridbAuth, ok = os.LookupEnv("SIRIDB_AUTH")
	if !ok {
		return errors.New("Missing SIRIDB_AUTHh in environment")
	}
	return
}

func DecodeAuth(in string) (av AuthVar, err error) {
	var buffer bytes.Buffer
	data, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error on decoding auth: %v", err))
		return
	}
	buffer.Write(data)
	bufstr := buffer.String()
	fmt.Printf("\n %s \n", bufstr)

	return
}
