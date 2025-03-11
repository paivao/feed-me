package utils

import (
	"database/sql/driver"
	"fmt"
	"net"
)

type MyNet struct {
	net.IPNet
}

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
	_, net, err := net.ParseCIDR(string(b))
	ip.IPNet = *net
	return err
}
