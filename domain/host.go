package core

import (
	"net"
)

type Host struct {
	Uid     int
	Name    string
	Address net.IP
	Comment string
}

func NewHost(name, addr, comment string) *Host {
	return &Host{
		Uid:     0,
		Name:    name,
		Address: net.ParseIP(addr),
		Comment: comment,
	}
}

func (h *Host) Match(addr string) bool {
	comp := net.ParseIP(addr)
	return h.Address.Equal(comp)
}
