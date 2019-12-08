package core

import (
	"fmt"
	"net"
)

type Host struct {
	Uid     int
	Name    string
	Address net.IP
	Comment string
}

func NewHost(name, addr, comment string) (*Host, error) {
	address := net.ParseIP(addr)
	if address == nil {
		return nil, fmt.Errorf("invalid host address: %s", addr)
	}
	return &Host{
		Uid:     0,
		Name:    name,
		Address: address,
		Comment: comment,
	}, nil
}

func (h *Host) Match(addr string) bool {
	comp := net.ParseIP(addr)
	return h.Address.Equal(comp)
}
