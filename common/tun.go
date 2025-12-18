package common

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/songgao/water"
)

type TunInterface struct {
	ifce *water.Interface
	name string
	ip   string
}

func NewTun(name, ip string) (*TunInterface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	ifce, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN interface: %v", err)
	}

	tun := &TunInterface{
		ifce: ifce,
		name: name,
		ip:   ip,
	}

	if err := tun.setupIP(); err != nil {
		ifce.Close()
		return nil, fmt.Errorf("failed to setup IP: %v", err)
	}

	if err := tun.setupLink(); err != nil {
		ifce.Close()
		return nil, fmt.Errorf("failed to setup link: %v", err)
	}

	return tun, nil
}

func (t *TunInterface) setupIP() error {
	cmd := exec.Command("ip", "addr", "add", t.ip, "dev", t.name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ip addr add failed: %v, output: %s", err, output)
	}
	return nil
}

func (t *TunInterface) setupLink() error {
	cmd := exec.Command("ip", "link", "set", t.name, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ip link set up failed: %v, output: %s", err, output)
	}
	return nil
}

func (t *TunInterface) Read(p []byte) (int, error) {
	return t.ifce.Read(p)
}

func (t *TunInterface) Write(p []byte) (int, error) {
	return t.ifce.Write(p)
}

func (t *TunInterface) Close() error {
	return t.ifce.Close()
}

func (t *TunInterface) Name() string {
	return t.name
}

func (t *TunInterface) IP() string {
	return t.ip
}

func AddRoute(route string) error {
	args := strings.Split(route, " ")
	fullArgs := append([]string{"route", "add"}, args...)
	cmd := exec.Command("ip", fullArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ip route add failed: %v, output: %s", err, output)
	}
	return nil
}

func GetDefaultGateway() (string, error) {
	cmd := exec.Command("ip", "route", "show", "default")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default route: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no default route found")
	}

	re := regexp.MustCompile(`via\s+(\S+)`)
	matches := re.FindStringSubmatch(lines[0])
	if matches == nil {
		return "", fmt.Errorf("could not parse gateway from route: %s", lines[0])
	}

	return matches[1], nil
}
