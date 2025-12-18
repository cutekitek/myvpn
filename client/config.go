package main

type Config struct {
	ServerAddr   string
	TunName      string
	TunIP        string
	ServerTunIP  string
	ServerRealIP string
	ClientID     uint64
}

func DefaultConfig() *Config {
	return &Config{
		ServerAddr:   ":5555",
		TunName:      "tun0",
		TunIP:        "10.0.0.1/24",
		ServerTunIP:  "10.0.0.2",
		ServerRealIP: "",
		ClientID:     1,
	}
}
