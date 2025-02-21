package mqtt

import (
	"crypto/tls"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttUrl string

type Mqtt struct {
	Host      string
	Port      string
	User      string
	Password  string
	Topics    []string
	CaCRT     []byte
	TlsConfig tls.Config
	Client    *mqtt.Client
}
