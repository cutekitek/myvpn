#!/bin/bash

set -e

echo "Setting up VPN client..."

SERVER_IP="$1"
if [ -z "$SERVER_IP" ]; then
    echo "Usage: $0 <server_ip>"
    echo "Example: $0 192.168.1.100"
    exit 1
fi

# Get default gateway
DEFAULT_GW=$(ip route | grep default | awk '{print $3}')
if [ -z "$DEFAULT_GW" ]; then
    echo "Error: Could not determine default gateway"
    exit 1
fi

echo "Server IP: $SERVER_IP"
echo "Default gateway: $DEFAULT_GW"

echo "Adding route for server via original gateway..."
ip route add $SERVER_IP via $DEFAULT_GW

echo "Note: Default route via TUN will be added by the client program"
echo ""
echo "To run the client:"
echo "  go run ./client -server $SERVER_IP:5555 -server-real-ip $SERVER_IP"
echo ""
echo "Make sure to run as root or with appropriate permissions."