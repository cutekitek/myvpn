package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	MaxPacketSize = 1500
	HeaderSize    = 8
)

type PacketHeader struct {
	ClientID uint32
	Length   uint32
}

func EncodePacket(clientID uint32, data []byte) []byte {
	buf := make([]byte, HeaderSize+len(data))
	binary.BigEndian.PutUint32(buf[0:4], clientID)
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(data)))
	copy(buf[8:], data)
	return buf
}

func DecodePacket(buf []byte) (uint32, []byte, error) {
	if len(buf) < HeaderSize {
		return 0, nil, fmt.Errorf("packet too small")
	}

	clientID := binary.BigEndian.Uint32(buf[0:4])
	length := binary.BigEndian.Uint32(buf[4:8])

	if uint32(len(buf)-HeaderSize) < length {
		return 0, nil, fmt.Errorf("packet truncated")
	}

	data := buf[HeaderSize : HeaderSize+int(length)]
	return clientID, data, nil
}

func IsIPv4(packet []byte) bool {
	if len(packet) < 1 {
		return false
	}
	version := packet[0] >> 4
	return version == 4
}

func GetDestinationIP(packet []byte) net.IP {
	if !IsIPv4(packet) || len(packet) < 20 {
		return nil
	}
	return net.IPv4(packet[16], packet[17], packet[18], packet[19])
}

func GetSourceIP(packet []byte) net.IP {
	if !IsIPv4(packet) || len(packet) < 20 {
		return nil
	}
	return net.IPv4(packet[12], packet[13], packet[14], packet[15])
}
