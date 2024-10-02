package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
)

type Net struct {
	net.IPNet
}

type IPEntry struct {
	Entry
	Network  net.IPNet `gorm:"type:binary(17);not null;uniqueIndex:idx_unique_entry"`
	IPFeedID uint      `gorm:"not null;uniqueIndex:idx_unique_entry"`
}

func (ip *Net) Scan(value interface{}) error {
	byte_arr, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to solve value: ", value))
	}
	var mask int = int(byte_arr[0])
	ip.IP = byte_arr[1:]
	ip.Mask = net.CIDRMask(mask, len(ip.IP)*8)
	return nil
}

// Value return binary value, implement driver.Valuer interface
func (ip Net) Value() (driver.Value, error) {
	ones, bits := ip.Mask.Size()
	if bits == 32 || bits == 128 {
		return nil, errors.New(fmt.Sprint("invalid IP length: ", bits))
	}
	bits = (bits >> 3) + 1
	data := make([]byte, bits)
	data[0] = byte(ones)
	copy(data[1:], ip.IP)
	return data, nil
}
