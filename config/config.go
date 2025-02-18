package config

/* Import libs */

/* Types */
type Config struct {
	ServerCfg struct {
		ListenIP string
		ListenPort uint
		MaxWorker uint
	}

	SiriDB struct {
		Host string
		Port uint
		Username string
		Password string
	}
	
	MongoDB struct {
		Url string
	}
	
	OpenWeather struct {
		Longitude float64
		Latitude float64
		Token string
	}

}

/* Global vars */

/* functions and methods */


// New -> create clean new instance of configuration
func New *Config {
	return &Config{}
}

// GetConfig -> read configuration, returns setuped instance of configuration struct
func GetConfig() (*Config, error) {
	cfg := new(Config)
	cfg.ServerCfg.ListenIP = "0.0.0.0"
	cfg.ServerCfg.ListenPort = 3500
	return cfg, nil
}

