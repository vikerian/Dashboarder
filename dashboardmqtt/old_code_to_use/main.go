package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	//	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/k0kubun/pp"
)

const MQTTHOST = "192.168.26.28"
const MQTTPORT = 8883
const MQTTUSER = "svc_zabbix"
const MQTTPASS = "GGyLTVRWEPlQP6wi6rwMNaZ55"

//const MQTTTOPIC = "#"

// const MQTTTOPIC = "/ttndata"
//const MQTTTOPIC = "/vodomery/#"

const MQTTTOPIC = "/#"

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//log.Printf("Received message for topic: %s\n", msg.Topic())
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	//log.Printf("%s", msg.Payload())
	go DecodePayload(msg)
}

type MessagePayload map[string]interface{}

/*type MessagePayload struct {
	CorrelationIDs interface{} `json:"correlation_ids,omitempty"`
	EndDeviceIDs   interface{} `json:"end_device-ids,omitempty"`
	ReceivedAT     []time.Time   `json:"received_at,omitempty"`
	UplinkMSG      struct {
		ConsumedAirTime time.Time `json:"consumed_airtime,omitempty"`
		DecodedPayload struct {
			Format string `json:"format"`
			Message []byte `json:"msg"`
			Sensor string `json:"sensor"`

		} `json:"decoded_payload"`
		F_CNT uint64 `json:"f_cnt"`
		F_PORT uint32 `json:"f_port"`
		FRM_PAYLOAD string `json:"frm_payload"`
		NetworkIDS struct {
			ClusterADDR string `json:"cluster_address"`
			ClusterID string `json:"cluster_id"`
			NetID uint64 `json:"net_id"`
			NsID []byte `json:"ns_id"`
			TenantID string `json:"tenant_id"`
		} `json:"network_ids"`
	} `json:"uplink_message,omitempty"`
}
*/
/*
type MessagePayload struct {
	CorrelationIds []string `json:"correlation_ids,omitempty"`
	EndDeviceIds   struct {
		ApplicationIds struct {
			ApplicationID string `json:"application_id,omitempty"`
		} `json:"application_ids,omitempty"`
		DevAddr  string `json:"dev_addr,omitempty"`
		DevEui   string `json:"dev_eui,omitempty"`
		DeviceID string `json:"device_id,omitempty"`
		JoinEui  string `json:"join_eui,omitempty"`
	} `json:"end_device_ids,omitempty"`
	ReceivedAt    time.Time `json:"received_at,omitempty"`
	UplinkMessage struct {
		ConsumedAirtime string `json:"consumed_airtime,omitempty"`
		DecodedPayload  struct {
			Format string `json:"format,omitempty"`
			Msg    string `json:"msg,omitempty"`
			Sensor string `json:"sensor,omitempty"`
		} `json:"decoded_payload,omitempty"`
		FCnt       int    `json:"f_cnt,omitempty"`
		FPort      int    `json:"f_port,omitempty"`
		FrmPayload string `json:"frm_payload,omitempty"`
		NetworkIds struct {
			ClusterAddress string `json:"cluster_address,omitempty"`
			ClusterID      string `json:"cluster_id,omitempty"`
			NetID          string `json:"net_id,omitempty"`
			NsID           string `json:"ns_id,omitempty"`
			TenantID       string `json:"tenant_id,omitempty"`
		} `json:"network_ids,omitempty"`
		ReceivedAt time.Time `json:"received_at,omitempty"`
		RxMetadata []struct {
			ChannelRssi int `json:"channel_rssi,omitempty"`
			GatewayIds  struct {
				Eui       string `json:"eui,omitempty"`
				GatewayID string `json:"gateway_id,omitempty"`
			} `json:"gateway_ids,omitempty"`
			Location struct {
				Altitude  int     `json:"altitude,omitempty"`
				Latitude  float64 `json:"latitude,omitempty"`
				Longitude float64 `json:"longitude,omitempty"`
				Source    string  `json:"source,omitempty"`
			} `json:"location,omitempty"`
			ReceivedAt  time.Time `json:"received_at,omitempty"`
			Rssi        int       `json:"rssi,omitempty"`
			Snr         int       `json:"snr,omitempty"`
			Time        time.Time `json:"time,omitempty"`
			Timestamp   int64     `json:"timestamp,2023/12/11 09:29:51 Received message: [{"misto": "zskamenicka_krcek", "stav": 111349}] from topic: /vodomery/decincoding_rate,omitempty"`
					SpreadingFactor int    `json:"spreading_factor,omitempty"`
				} `json:"lora,omitempty"`
			} `json:"data_rate,omitempty"`
			Frequency string    `json:"frequency,omitempty"`
			Time      time.Time `json:"time,omitempty"`
			Timestamp int64     `json:"timestamp,omitempty"`
		} `json:"settings,omitempty"`
	} `json:"uplink_message,omitempty"`
}
*/

func DecodePayload(msg mqtt.Message) (data interface{}, err error) {
	rawData := msg.Payload()
	decData := json.NewDecoder(bytes.NewReader(rawData))
	var m MessagePayload
	err = decData.Decode(&m)
	data = m
	pp.Print(m)
	//fmt.Printf(m)
	return
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Printf("Connected\n")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v \n", err)
}

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.pem")
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

func sub(client mqtt.Client) {
	//topic := "/ttndata"

	//topic := "#"
	topic := MQTTTOPIC
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic %s\n", topic)
}

func main() {
	log.Printf("Starting up...\n")
	defer log.Printf("Ending...\n")
	var broker = MQTTHOST
	var port = MQTTPORT
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	tlsConfig := NewTlsConfig()
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(MQTTUSER)
	opts.SetPassword(MQTTPASS)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	fmt.Printf("mqtt client: n")
	pp.Print(client)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Printf("MQTT token: %+v", token)

	sub(client)
	for {

	}

}
