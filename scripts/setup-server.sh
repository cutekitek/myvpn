#!/bin/bash

set -e

echo "Setting up VPN server..."

# Enable IP forwarding
echo "Enabling IP forwarding..."
sysctl -w net.ipv4.ip_forward=1
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf

# Set up NAT masquerade for client subnet
echo "Setting up NAT masquerade..."
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE

# Allow forwarding between interfaces
echo "Setting up forwarding rules..."
iptables -A FORWARD -i tun0 -o eth0 -j ACCEPT
iptables -A FORWARD -i eth0 -o tun0 -m state --state ESTABLISHED,RELATED -j ACCEPT

# Save iptables rules
echo "Saving iptables rules..."
if command -v iptables-save &> /dev/null; then
    iptables-save > /etc/iptables.rules
fi

echo "Server setup complete!"
echo ""
echo "To run the server:"
echo "  go run ./server -port 5555"
echo ""
echo "Make sure to run as root or with appropriate permissions."