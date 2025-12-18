# Simple UDP Tunnel Proxy

A Go-based UDP tunnel proxy that routes all system traffic through a TUN interface, encapsulates it in UDP packets, and forwards to a remote server.

## Architecture

```
Client System → Client TUN → UDP encapsulation → Server UDP → Server TUN → Internet
Internet → Server TUN → UDP encapsulation → Client UDP → Client TUN → Client System
```

## Prerequisites

- Go 1.16+
- Root/sudo privileges for TUN interface creation
- Linux with TUN/TAP support

## Installation

```bash
git clone <repository>
cd vpns
go mod download
```

## Server Setup

1. Run setup script (as root):
```bash
sudo ./scripts/setup-server.sh
```

2. Build and run server:
```bash
cd server
go build
sudo ./server -port 5555
```

## Client Setup

1. Run setup script (as root):
```bash
sudo ./scripts/setup-client.sh <server_ip>
```

2. Build and run client:
```bash
cd client
go build
sudo ./client -server <server_ip>:5555 -server-real-ip <server_ip>
```

## Configuration

### Server Flags
- `-port`: UDP port to listen on (default: 5555)
- `-tun`: TUN interface name (default: tun0)
- `-tun-ip`: TUN interface IP (default: 10.0.0.2/24)
- `-client-subnet`: Client subnet (default: 10.0.0.0/24)

### Client Flags
- `-server`: Server address (host:port)
- `-tun`: TUN interface name (default: tun0)
- `-tun-ip`: TUN interface IP (default: 10.0.0.1/24)
- `-server-tun-ip`: Server TUN IP (default: 10.0.0.2)
- `-server-real-ip`: Server real IP address (required)
- `-client-id`: Client ID (default: 1)

## Network Configuration

### Server IP Forwarding
The server requires IP forwarding and NAT rules:
```bash
sysctl net.ipv4.ip_forward=1
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE
iptables -A FORWARD -i tun0 -o eth0 -j ACCEPT
iptables -A FORWARD -i eth0 -o tun0 -m state --state ESTABLISHED,RELATED -j ACCEPT
```

### Client Routing
Client adds routes:
```bash
# Route server IP via original gateway
ip route add <server_ip> via <original_gateway>

# Route all other traffic via TUN
ip route add default via 10.0.0.2 dev tun0
```

## Testing

1. Start server on remote VPS
2. Start client on local machine
3. Test connectivity:
```bash
ping 8.8.8.8
curl https://ifconfig.me
```

## Notes

- No encryption or authentication (as per requirements)
- Simple UDP encapsulation
- IPv4 traffic only
- Basic error handling and reconnection logic
- Requires root privileges for interface configuration

## Troubleshooting

1. **Permission denied**: Run with sudo
2. **TUN device creation failed**: Ensure TUN/TAP module is loaded
3. **No internet connectivity**: Check server IP forwarding and iptables rules
4. **Cannot reach server**: Verify routes and firewall rules