package mqtt

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
	config *Config
}

type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

// NewClient - create new client instance
func NewClient(config *Config) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)
	if config.Username != "" {
		opts.SetUsername(config.Username)
		opts.SetPassword(config.Password)
	}
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{
		client: client,
		config: config,
	}, nil
}

// Publish - publish data / command into topic
func (c *Client) Publish(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	token := c.client.Publish(topic, 1, false, data)
	token.Wait()
	return token.Error()
}

// Subscribe to topic for data
func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, 1, handler)
	token.Wait()
	return token.Error()
}

// Disconnect from mqtt
func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}
