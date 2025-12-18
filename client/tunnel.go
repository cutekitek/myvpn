package main

import (
	"fmt"
	"log"
	"vpns/common"
)

type Tunnel struct {
	config  *Config
	tun     *common.TunInterface
	udp     *common.UDPClient
	stop    chan struct{}
	running bool
}

func NewTunnel(config *Config) (*Tunnel, error) {
	tun, err := common.NewTun(config.TunName, config.TunIP)
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN interface: %v", err)
	}

	udp, err := common.NewUDPClient(config.ServerAddr)
	if err != nil {
		tun.Close()
		return nil, fmt.Errorf("failed to create UDP client: %v", err)
	}

	return &Tunnel{
		config: config,
		tun:    tun,
		udp:    udp,
		stop:   make(chan struct{}),
	}, nil
}

func (t *Tunnel) Start() error {
	if t.running {
		return fmt.Errorf("tunnel already running")
	}

	t.running = true

	if err := t.setupRouting(); err != nil {
		return fmt.Errorf("failed to setup routing: %v", err)
	}

	go t.readFromTun()
	go t.readFromUDP()

	log.Printf("Tunnel started: TUN=%s, Server=%s", t.tun.Name(), t.config.ServerAddr)
	return nil
}

func (t *Tunnel) Stop() {
	if !t.running {
		return
	}

	close(t.stop)
	t.running = false
	t.udp.Close()
	t.tun.Close()

	log.Println("Tunnel stopped")
}

func (t *Tunnel) setupRouting() error {
	if t.config.ServerRealIP == "" {
		return fmt.Errorf("server real IP is required for routing")
	}

	defaultGW, err := common.GetDefaultGateway()
	if err != nil {
		log.Printf("Warning: Could not get default gateway: %v", err)
	} else {
		serverRoute := fmt.Sprintf("%s via %s", t.config.ServerRealIP, defaultGW)
		if err := common.AddRoute(serverRoute); err != nil {
			log.Printf("Warning: Failed to add server route: %v", err)
		}
	}

	defaultRoute := fmt.Sprintf("default via %s dev %s", t.config.ServerTunIP, t.config.TunName)
	if err := common.AddRoute(defaultRoute); err != nil {
		return fmt.Errorf("failed to add default route: %v", err)
	}

	return nil
}

func (t *Tunnel) readFromTun() {
	buf := make([]byte, common.MaxPacketSize)

	for {
		select {
		case <-t.stop:
			return
		default:
			n, err := t.tun.Read(buf)
			if err != nil {
				log.Printf("Error reading from TUN: %v", err)
				continue
			}

			packet := buf[:n]
			encoded := common.EncodePacket(uint32(t.config.ClientID), packet)

			if err := t.udp.Write(encoded); err != nil {
				log.Printf("Error writing to UDP: %v", err)
			}
		}
	}
}

func (t *Tunnel) readFromUDP() {
	buf := make([]byte, common.MaxPacketSize+common.HeaderSize)

	for {
		select {
		case <-t.stop:
			return
		default:
			n, err := t.udp.Read(buf)
			if err != nil {
				log.Printf("Error reading from UDP: %v", err)
				continue
			}

			clientID, packet, err := common.DecodePacket(buf[:n])
			if err != nil {
				log.Printf("Error decoding packet: %v", err)
				continue
			}

			if uint64(clientID) != t.config.ClientID {
				log.Printf("Received packet for wrong client ID: %d", clientID)
				continue
			}

			if _, err := t.tun.Write(packet); err != nil {
				log.Printf("Error writing to TUN: %v", err)
			}
		}
	}
}
