package common

import (
	"fmt"
	"net"
)

type UDPServer struct {
	conn *net.UDPConn
	port int
}

func NewUDPServer(port int) (*UDPServer, error) {
	addr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on UDP port %d: %v", port, err)
	}

	return &UDPServer{
		conn: conn,
		port: port,
	}, nil
}

func (s *UDPServer) ReadFrom(buf []byte) (int, *net.UDPAddr, error) {
	return s.conn.ReadFromUDP(buf)
}

func (s *UDPServer) WriteTo(buf []byte, addr *net.UDPAddr) error {
	_, err := s.conn.WriteToUDP(buf, addr)
	return err
}

func (s *UDPServer) Close() error {
	return s.conn.Close()
}

func (s *UDPServer) LocalAddr() *net.UDPAddr {
	return s.conn.LocalAddr().(*net.UDPAddr)
}

type UDPClient struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func NewUDPClient(serverAddr string) (*UDPClient, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %v", err)
	}

	return &UDPClient{
		conn: conn,
		addr: addr,
	}, nil
}

func (c *UDPClient) Read(buf []byte) (int, error) {
	return c.conn.Read(buf)
}

func (c *UDPClient) Write(buf []byte) error {
	_, err := c.conn.Write(buf)
	return err
}

func (c *UDPClient) Close() error {
	return c.conn.Close()
}

func (c *UDPClient) RemoteAddr() *net.UDPAddr {
	return c.addr
}
