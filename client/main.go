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

	flag.StringVar(&config.ServerAddr, "server", config.ServerAddr, "Server address (host:port)")
	flag.StringVar(&config.TunName, "tun", config.TunName, "TUN interface name")
	flag.StringVar(&config.TunIP, "tun-ip", config.TunIP, "TUN interface IP address")
	flag.StringVar(&config.ServerTunIP, "server-tun-ip", config.ServerTunIP, "Server TUN IP address")
	flag.StringVar(&config.ServerRealIP, "server-real-ip", config.ServerRealIP, "Server real IP address")
	flag.Uint64Var(&config.ClientID, "client-id", config.ClientID, "Client ID")
	flag.IntVar(&config.MTU, "mtu", config.MTU, "MTU for TUN interface (0 = default 1464)")
	flag.Parse()

	if config.ServerRealIP == "" {
		log.Fatal("Server real IP is required (use -server-real-ip flag)")
	}

	tunnel, err := NewTunnel(config)
	if err != nil {
		log.Fatalf("Failed to create tunnel: %v", err)
	}

	if err := tunnel.Start(); err != nil {
		log.Fatalf("Failed to start tunnel: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("VPN client started. Press Ctrl+C to stop.")
	<-sigChan

	tunnel.Stop()
	log.Println("VPN client stopped.")
}
