package mqtt

import (
	"crypto/tls"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttUrl string

type Mqtt struct {
	Url       string // Url is tcp://user:password@host:port
	ClientID  string // Set on Newclient, otherwise generated uuid -> last 8 chars without --
	opts	  *mqtt.ClientOptions
	Topics    []string // 
	caCrtPATH string
	CaCRT     []byte
	TlsConfig tls.Config
	Client    mqtt.Client
}

/* NewClient => constructor
 *	-> InputParams: 
 *		string uri -> url in format tcp://user:pass@host:port
 *		string caCrtPath -> path on filesystem to CA certificate (yes have to be mounted RO to container)
 *	-> return clean instance filled with params from constructor
*/
func NewClient(string uri, string caCrtPath) *Mqtt {
	return &Mqtt{
		Url: uri,
		caCrtPATH: caCrtPath
	}
}

/* SetupConnection -> setup parameters for connection
 * inputs:
 *   string username
 *   string password
 *   if tls, then *tls.Config, otherwise nil
 * returns:
 *   bool ok
 *   error if any occured
*/
func (mqt *Mqtt) SetupConnection(user string, password string, tlsConf *tls.Config, ) (bool ok, err error) {
	
	return
}	
