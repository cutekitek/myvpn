package main

type Config struct {
	ListenPort   int
	TunName      string
	TunIP        string
	ClientSubnet string
}

func DefaultConfig() *Config {
	return &Config{
		ListenPort:   5555,
		TunName:      "tun0",
		TunIP:        "10.0.0.2/24",
		ClientSubnet: "10.0.0.0/24",
	}
}
