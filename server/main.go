package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := DefaultConfig()

	flag.IntVar(&config.ListenPort, "port", config.ListenPort, "UDP port to listen on")
	flag.StringVar(&config.TunName, "tun", config.TunName, "TUN interface name")
	flag.StringVar(&config.TunIP, "tun-ip", config.TunIP, "TUN interface IP address")
	flag.StringVar(&config.ClientSubnet, "client-subnet", config.ClientSubnet, "Client subnet")
	flag.IntVar(&config.MTU, "mtu", config.MTU, "MTU for TUN interface (0 = default 1464)")
	flag.Parse()

	server, err := NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("VPN server started. Press Ctrl+C to stop.")
	<-sigChan

	server.Stop()
	log.Println("VPN server stopped.")
}
