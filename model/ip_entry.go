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
	Network  net.IPNet `gorm:"type:binary(17)"`
	IPFeedID uint
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

// Value return json value, implement driver.Valuer interface
func (ip Net) Value() (driver.Value, error) {
	return ip.String(), nil
}
