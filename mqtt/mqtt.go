package mqtt

import (
	"crypto/tls"
//	"crypto/x509"
	"fmt"
//	"io/ioutil"
//	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/google/uuid"
)

const ClientPrefix string = "dcukSu"
const CaCrtPath string = "./ca.pem"

var mqttUrl string

type Mqtt struct {
	URL       string // Url is tcp://user:password@host:port
	ClientID  string // Set on Newclient, otherwise generated uuid -> last 8 chars without --
	Opts      *mqtt.ClientOptions
	Topics    []string //
	caCrtPATH string
	CaCRT     []byte
	TlsCONFIG *tls.Config
	CLH    mqtt.Client
}

/* New => constructor
 *	-> InputParams:
 *		string uri -> url in format tcp://user:pass@host:port
 *		string caCrtPath -> path on filesystem to CA certificate (yes have to be mounted RO to container)
 *	-> return clean instance filled with params from constructor
 */
func New(uri, caCrtPath string) (mq Mqtt,err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return 
	}
	clientID := fmt.Sprintf("%s-%s", ClientPrefix, id.String()[4:8])
	mq.URL = uri
	mq.ClientID = clientID
	mq.caCrtPATH = caCrtPath
	mq.Opts = mqtt.NewClientOptions().AddBroker(mq.URL).SetClientID(clientID)
	mq.CLH = mqtt.NewClient(mq.Opts)
	if token := mq.CLH.Connect(); token.Wait() && token.Error() != nil {
		return
	}
	return
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

*/
