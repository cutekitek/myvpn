# Simple UDP Tunnel Proxy - Implementation Plan

## Overview
Create a Go-based UDP tunnel proxy that routes all system traffic through a TUN interface, encapsulates it in UDP packets, sends to a remote server which forwards traffic to the Internet, and returns responses back to the client. Uses the `github.com/songgao/water` library for TUN interface management.

## Architecture

### Components

**1. Client** - Runs on user machine
- Creates TUN interface with assigned IP (e.g., `10.0.0.1/24`)
- Sets up routing to send traffic through TUN
- Reads packets from TUN, encapsulates in UDP, sends to server
- Receives UDP from server, writes to TUN

**2. Server** - Runs on remote VPS
- Listens on UDP port (e.g., `5555`)
- Creates TUN interface with IP (e.g., `10.0.0.2/24`)
- Maintains mapping between client TUN IPs and UDP addresses
- Forwards packets between TUN and UDP clients
- Acts as gateway for Internet traffic

### Packet Flow
```
Client System → Client TUN → UDP encapsulation → Server UDP → Server TUN → Internet
Internet → Server TUN → UDP encapsulation → Client UDP → Client TUN → Client System
```

## Technical Design

### Client Implementation
- Use `water.New(water.Config{DeviceType: water.TUN})` to create interface
- Assign IP address using `ip` commands or netlink
- Set up routing: default route via TUN, exclude server IP
- Two goroutines:
  1. TUN → UDP: read packets, send to server
  2. UDP → TUN: receive packets, write to TUN

### Server Implementation
- UDP listener on configurable port
- TUN interface for network access
- IP forwarding enabled (`net.ipv4.ip_forward=1`)
- NAT masquerade for client subnet
- Concurrent client support with IP→UDP address mapping
- Two goroutines per client:
  1. UDP → TUN: forward client traffic to Internet
  2. TUN → UDP: forward Internet responses to client

### IP Addressing
- Client TUN: `10.0.0.1/24`
- Server TUN: `10.0.0.2/24`
- UDP port: `5555` (configurable)

## Required System Configuration

### Server Configuration
```bash
# Enable IP forwarding
sysctl net.ipv4.ip_forward=1

# NAT masquerade for client subnet
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE

# Allow forwarding between interfaces
iptables -A FORWARD -i tun0 -o eth0 -j ACCEPT
iptables -A FORWARD -i eth0 -o tun0 -m state --state ESTABLISHED,RELATED -j ACCEPT
```

### Client Configuration
```bash
# Route server IP via original gateway
ip route add <server_ip> via <original_gateway>

# Route all other traffic via TUN
ip route add default via 10.0.0.2 dev tun0
```

## Dependencies
- Go 1.16+
- `github.com/songgao/water` for TUN interfaces
- `github.com/songgao/water/waterutil` for packet parsing
- Root privileges for interface configuration

## Implementation Steps

1. **Setup Project Structure**
   - Create Go module
   - Install dependencies
   - Create directory structure

2. **Core Library Components**
   - TUN interface wrapper
   - UDP tunnel implementation
   - IP address management

3. **Client Implementation**
   - Command-line interface
   - Configuration parsing
   - Routing setup automation

4. **Server Implementation**
   - Multi-client UDP server
   - Connection state management
   - NAT and forwarding logic

5. **Testing and Validation**
   - Unit tests for core components
   - Integration testing with real traffic
   - Performance benchmarking

6. **Documentation and Deployment**
   - Usage instructions
   - Installation scripts
   - Systemd service files

## File Structure
```
vpns/
├── go.mod
├── go.sum
├── client/
│   ├── main.go
│   ├── config.go
│   └── tunnel.go
├── server/
│   ├── main.go
│   ├── config.go
│   └── tunnel.go
├── common/
│   ├── tun.go
│   ├── packet.go
│   └── udp.go
└── scripts/
    ├── setup-server.sh
    └── setup-client.sh
```

## Notes
- No encryption or authentication as per requirements
- Simple UDP encapsulation without complex framing
- Basic error handling and reconnection logic
- Support for IPv4 traffic only initially
- MTU considerations for UDP fragmentation