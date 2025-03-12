package utils

import (
	"database/sql/driver"
	"fmt"
	"net"
	"strings"
)

type MyNet struct {
	net.IPNet
}

var ipv4Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255}

func (ip *MyNet) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to parse IP/Net value: %v", src)
	}
	var mask int = int(bytes[0])
	bytes = bytes[1:]
	if len(bytes) != 4 && len(bytes) != 16 {
		return fmt.Errorf("invalid IP length: %d", len(bytes))
	}
	ip.IP = make(net.IP, len(bytes))
	copy(ip.IP, bytes)
	ip.Mask = net.CIDRMask(mask, len(ip.IP)*8)
	return nil
}

func (ip MyNet) Value() (driver.Value, error) {
	ones, bits := ip.Mask.Size()
	if bits != 32 && bits != 128 {
		return nil, fmt.Errorf("invalid IP length: %d", bits)
	}
	bits = (bits >> 3) + 1
	data := make([]byte, bits)
	data[0] = byte(ones)
	copy(data[1:], ip.IP)
	return data, nil
}

func (ip *MyNet) MarshalJSON() ([]byte, error) {
	return []byte(ip.String()), nil
}

func (ip *MyNet) UnmarshalJSON(b []byte) error {
	value := string(b[1 : len(b)-1])
	if strings.IndexByte(value, '/') == -1 {
		calculated_ip := net.ParseIP(value)
		calculated_ipv4 := calculated_ip.To4()
		if calculated_ip == nil {
			return &net.ParseError{Type: "Unknown IP address", Text: value}
		}
		if calculated_ipv4 != nil {
			ip.IP = calculated_ipv4
			ip.Mask = net.CIDRMask(net.IPv4len*8, net.IPv4len*8)
		} else {
			ip.IP = calculated_ip
			ip.Mask = net.CIDRMask(net.IPv6len*8, net.IPv6len*8)
		}
		return nil
	}
	_, net, err := net.ParseCIDR(value)
	if err != nil {
		return err
	}
	ip.IPNet = *net
	return nil
}
