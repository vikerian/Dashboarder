package mqtt

import (
	"crypto/tls"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttUrl string

type Mqtt struct {
	Log       *slog.Logger
	Host      string
	Port      string
	User      string
	Password  string
	Topics    []string
	CaCRT     []byte
	TlsConfig tls.Config
	Client    *mqtt.Client
}

// New -> return clean instance
func New() *Mqtt {
	return new(Mqtt)
}

// Setup -> setup instance params
func (mqt *Mqtt) SetupParam(key string, val interface{}) {
	switch key {
	case "log":
		mqt.Log = val.(*slog.Logger)
	case "host":
		mqt.Host = val.(string)
	case "port":
		mqt.Port = val.(string)
	case "user":
		mqt.User = val.(string)
	case "username":
		mqt.User = val.(string)
	case "pass":
		mqt.Password = val.(string)
	case "password":
		mqt.Password = val.(string)
	case "topics":
		mqt.Topics = val.([]string)
	case "ca-crt":
		mqt.CaCRT = val.([]byte)
	case "cacrt":
		mqt.CaCRT = val.([]byte)
	case "ca_crt":
		mqt.CaCRT = val.([]byte)
	case "tlsconfig":
		mqt.TlsConfig = val.(tls.Config)
	case "tls-config":
		mqt.TlsConfig = val.(tls.Config)
	case "tls_config":
		mqt.TlsConfig = val.(tls.Config)
	}
}

// Connect -> connect to mqtt and authorize ourselves
func (mqt *Mqtt) Connect()
