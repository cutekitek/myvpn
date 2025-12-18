package common

import (
	"testing"
)

func TestEncodeDecodePacket(t *testing.T) {
	clientID := uint32(12345)
	data := []byte("test packet data")

	encoded := EncodePacket(clientID, data)
	decodedID, decodedData, err := DecodePacket(encoded)

	if err != nil {
		t.Fatalf("DecodePacket failed: %v", err)
	}

	if decodedID != clientID {
		t.Errorf("Expected client ID %d, got %d", clientID, decodedID)
	}

	if string(decodedData) != string(data) {
		t.Errorf("Expected data %s, got %s", data, decodedData)
	}
}

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		name   string
		packet []byte
		want   bool
	}{
		{"IPv4 packet", []byte{0x45, 0x00}, true},
		{"IPv6 packet", []byte{0x60, 0x00}, false},
		{"Empty packet", []byte{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIPv4(tt.packet)
			if got != tt.want {
				t.Errorf("IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDestinationIP(t *testing.T) {
	packet := []byte{
		0x45, 0x00, 0x00, 0x14,
		0x00, 0x00, 0x00, 0x00,
		0x40, 0x00, 0x00, 0x00,
		0xc0, 0xa8, 0x00, 0x01, // Source: 192.168.0.1
		0x08, 0x08, 0x08, 0x08, // Destination: 8.8.8.8
	}

	ip := GetDestinationIP(packet)
	if ip == nil {
		t.Fatal("GetDestinationIP returned nil")
	}

	if ip.String() != "8.8.8.8" {
		t.Errorf("Expected destination IP 8.8.8.8, got %s", ip.String())
	}
}
