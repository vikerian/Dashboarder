package main

type Config struct {
	MqttURL         string
	MqttAUTH        string
	MqttChannelLIST []string
	ValkeyURL       string
	ValkeyAUTH      string
	SiriDbURL       string
	SiriDbAUTH      string
}

func NewConfig() Config {
	cf := new(Config)
	return *cf
}

func (cf *Config) LoadConfiguration() (ok bool, err error) {

	return
}
