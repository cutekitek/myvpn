package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"vpns/common"
)

type Client struct {
	ID       uint32
	Addr     *net.UDPAddr
	LastSeen int64
}

type Server struct {
	config     *Config
	tun        *common.TunInterface
	udp        *common.UDPServer
	clients    map[uint32]*Client
	ipToClient map[string]*Client
	mu         sync.RWMutex
	stop       chan struct{}
	running    bool
}

func NewServer(config *Config) (*Server, error) {
	tun, err := common.NewTun(config.TunName, config.TunIP)
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN interface: %v", err)
	}

	udp, err := common.NewUDPServer(config.ListenPort)
	if err != nil {
		tun.Close()
		return nil, fmt.Errorf("failed to create UDP server: %v", err)
	}

	return &Server{
		config:     config,
		tun:        tun,
		udp:        udp,
		clients:    make(map[uint32]*Client),
		ipToClient: make(map[string]*Client),
		stop:       make(chan struct{}),
	}, nil
}

func (s *Server) Start() error {
	if s.running {
		return fmt.Errorf("server already running")
	}

	s.running = true

	if err := s.setupForwarding(); err != nil {
		return fmt.Errorf("failed to setup forwarding: %v", err)
	}

	go s.handleUDPClients()
	go s.readFromTun()

	log.Printf("Server started on UDP port %d, TUN=%s", s.config.ListenPort, s.tun.Name())
	return nil
}

func (s *Server) Stop() {
	if !s.running {
		return
	}

	close(s.stop)
	s.running = false
	s.udp.Close()
	s.tun.Close()

	log.Println("Server stopped")
}

func (s *Server) setupForwarding() error {
	log.Println("Note: Ensure IP forwarding is enabled on the system:")
	log.Println("  sysctl net.ipv4.ip_forward=1")
	log.Println("  iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE")
	log.Println("  iptables -A FORWARD -i tun0 -o eth0 -j ACCEPT")
	log.Println("  iptables -A FORWARD -i eth0 -o tun0 -m state --state ESTABLISHED,RELATED -j ACCEPT")
	return nil
}

func (s *Server) handleUDPClients() {
	buf := make([]byte, common.MaxPacketSize+common.HeaderSize)

	for {
		select {
		case <-s.stop:
			return
		default:
			n, addr, err := s.udp.ReadFrom(buf)
			if err != nil {
				log.Printf("Error reading from UDP: %v", err)
				continue
			}

			clientID, packet, err := common.DecodePacket(buf[:n])
			if err != nil {
				log.Printf("Error decoding packet: %v", err)
				continue
			}

			s.mu.Lock()
			client, exists := s.clients[clientID]
			if !exists {
				client = &Client{
					ID:   clientID,
					Addr: addr,
				}
				s.clients[clientID] = client
				log.Printf("New client connected: ID=%d, Addr=%s", clientID, addr)
			} else {
				client.Addr = addr
			}
			if sourceIP := common.GetSourceIP(packet); sourceIP != nil {
				s.ipToClient[sourceIP.String()] = client
			}
			s.mu.Unlock()

			if _, err := s.tun.Write(packet); err != nil {
				log.Printf("Error writing to TUN: %v", err)
			}
		}
	}
}

func (s *Server) readFromTun() {
	buf := make([]byte, common.MaxPacketSize+4)

	for {
		select {
		case <-s.stop:
			return
		default:
			n, err := s.tun.Read(buf[4:])
			if err != nil {
				log.Printf("Error reading from TUN: %v", err)
				continue
			}

			packet := buf[4:n]
			destIP := common.GetDestinationIP(packet)
			if destIP != nil {
				s.mu.RLock()
				client, found := s.ipToClient[destIP.String()]
				s.mu.RUnlock()
				if found {
					binary.BigEndian.PutUint32(buf, client.ID)
					if err := s.udp.WriteTo(buf[:n], client.Addr); err != nil {
						log.Printf("Error sending to client %d: %v", client.ID, err)
					}
				} else {
					log.Printf("No client found for destination IP: %s", destIP)
				}
			}
		}
	}
}

func (s *Server) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}
