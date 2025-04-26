package main

import "os"

type Config struct {
	MqttURL         string   `json:"mqtt_url"`
	MqttAUTH        string   `json:"mqtt_auth"`
	MqttChannelLIST []string `json:"mqtt_channels"`
	ValkeyURL       string   `json:"valkey_url"`
	ValkeyAUTH      string   `json:"valkey_auth"`
	SiriDbURL       string   `json:"siridb_url"`
	SiriDbAUTH      string   `json:"siridb_auth"`
}

func NewConfig() Config {
	cf := new(Config)
	return *cf
}

func (cf *Config) LoadConfiguration() (ok bool, err error) {
	cf.MqttURL,ok = os.Lookupenv("MQTT_URL")
	cf.MqttAUTH,ok = os.
	return
}
