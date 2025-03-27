package types

import (
	"bytes"
	"errors"
	"net/netip"
)

type MyNet struct {
	netip.Prefix
}

func (ip *MyNet) ScanBytes(src []byte) error {
	return ip.UnmarshalBinary(src)
}

// Must not be a pointer
func (ip MyNet) BytesValue() ([]byte, error) {
	return ip.MarshalBinary()
}

func (ip MyNet) MarshalJSON() ([]byte, error) {
	if !ip.IsValid() {
		return nil, errors.New("invalid IP/Net")
	}
	val := ip.ToString()
	b := make([]byte, 0, len(val)+2)
	b = append(b, '"')
	b = append(b, val...)
	b = append(b, '"')
	return b, nil
}

func (ip *MyNet) UnmarshalJSON(b []byte) error {
	value := b[1 : len(b)-1]
	return ip.FromText(value)
}

func (ip *MyNet) FromText(value []byte) error {
	if bytes.IndexByte(value, '/') == -1 {
		//ipv4
		if bytes.IndexByte(value, ':') == -1 {
			value = append(value, "/32"...)
		} else {
			value = append(value, "/128"...)
		}
	}
	return ip.UnmarshalText(value)
}

func (ip *MyNet) ToString() string {
	var val string
	if ip.IsSingleIP() {
		val = ip.Addr().String()
	} else {
		val = ip.String()
	}
	return val
}

func (ip *MyNet) String() string {
	return ip.ToString()
}
