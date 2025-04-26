package main

// import libs
import (
	"encoding/base64"
	"errors"
	"os"
)

// types
type Config struct {
	MqttUrl         string   `json:"mqtt_url"`
	MqttAuth        string   `json:"mqtt_auth"`
	MqttChannelLIST []string `json:"mqtt_channels"`
	ValkeyUrl       string   `json:"valkey_url"`
	ValkeyAuth      string   `json:"valkey_auth"`
	SiridbUrl       string   `json:"siridb_url"`
	SiridbAuth      string   `json:"siridb_auth"`
}

type AuthVar struct {
}

// constructors
func NewConfig() Config {
	cf := new(Config)
	return *cf
}

// methods
func (cf *Config) LoadConfiguration() (ok bool, err error) {
	cf.MqttUrl, ok = os.LookupEnv("MQTT_Url")
	if !ok {
		return false, errors.New("Missing MQTT_Url in environment")
	}
	cf.MqttAuth, ok = os.LookupEnv("MQTT_Auth")
	if !ok {
		return false, errors.New("Missing MQTT_Auth in environment")
	}
	cf.ValkeyUrl, ok = os.LookupEnv("VALKEY_Url")
	if !ok {
		return false, errors.New("Missing VALKEY_Url in environment")
	}
	cf.ValkeyAuth, ok = os.LookupEnv("VALKEY_Auth")
	if !ok {
		return false, errors.New("Missing VALKEY_Auth in environment")
	}
	cf.SiridbUrl, ok = os.LookupEnv("SIRIdb_Url")
	if !ok {
		return false, errors.New("Missing SIRIdb_Url in environment")
	}
	cf.SiridbAuth, ok = os.LookupEnv("SIRIdb_Auth")
	if !ok {
		return false, errors.New("Missing SIRIdb_Auth in environment")
	}
	return
}

func (cf *Config) DecodeAuth(instr string) (av Authvar, err error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(intstr)
	if err != nil {
		return
	}
	return
}
