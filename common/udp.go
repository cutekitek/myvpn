package common

import (
	"fmt"
	"log"
	"net"
)

func setSocketBuffers(conn *net.UDPConn) {
	if err := conn.SetReadBuffer(SocketBufferSize); err != nil {
		log.Printf("Warning: failed to set UDP read buffer size: %v", err)
	}
	if err := conn.SetWriteBuffer(SocketBufferSize); err != nil {
		log.Printf("Warning: failed to set UDP write buffer size: %v", err)
	}
}

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

	setSocketBuffers(conn)

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

	setSocketBuffers(conn)

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
