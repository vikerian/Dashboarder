package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/google/uuid"
)

const ClientPrefix string = "dcukSu"
const CaCrtPath string = "./ca.pem"

var mqttUrl string

type Mqtt struct {
	Url       string // Url is tcp://user:password@host:port
	ClientID  string // Set on Newclient, otherwise generated uuid -> last 8 chars without --
	opts      *mqtt.ClientOptions
	Topics    []string //
	caCrtPATH string
	CaCRT     []byte
	TlsConfig *tls.Config
	Client    mqtt.Client
}

/* NewClient => constructor
 *	-> InputParams:
 *		string uri -> url in format tcp://user:pass@host:port
 *		string caCrtPath -> path on filesystem to CA certificate (yes have to be mounted RO to container)
 *	-> return clean instance filled with params from constructor
 */
func NewClient(uri, caCrtPath string) (*Mqtt, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	clientID := fmt.Sprintf("%s-%s", ClientPrefix, id.String()[4:8])
	return &Mqtt{
		Url:       uri,
		ClientID:  clientID,
		caCrtPATH: caCrtPath,
	}, nil
}

/* Connect -> connect to mqtt and get version info
 * 

func newTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(CaCrtPath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	}
}
