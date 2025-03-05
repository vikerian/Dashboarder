package mqtt

import (
	"crypto/tls"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttUrl string

type Mqtt struct {
	Url       string
	ClientID  string
	Topics    []string
	CaCRT     []byte
	TlsConfig tls.Config
	Client    mqtt.Client
}

// New -> return clean instance
func New() *Mqtt {
	return new(Mqtt)
}

// Setup -> setup instance params
func (mqt *Mqtt) SetupParam(key string, val interface{}) {
	switch key {
	case "url", "URL", "Url":
		mqt.Url = val.(string)
	case "topics":
		mqt.Topics = val.([]string)
	case "ca-crt", "cacrt", "ca_crt":
		mqt.CaCRT = val.([]byte)
	case "tlsconfig", "tls-config", "tls_config":
		mqt.TlsConfig = val.(tls.Config)
	case "clientID", "client_id":
		mqt.ClientID = val.(string)
	}
}

// Connect -> connect to mqtt and authorize ourselves
func (mqt *Mqtt) Connect() (err error) {
	// create client with options
	clhOPTS := mqtt.NewClientOptions()
	// add login url
	clhOPTS.AddBroker(mqt.Url)
	// add client id
	clhOPTS.SetClientID(mqt.ClientID)
	mqt.Client = mqtt.NewClient(clhOPTS)
	if token := mqt.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
