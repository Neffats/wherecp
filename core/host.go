package core

import (
	"fmt"
	"net"
)

// Host represents a single IPv4 host object, a single IPv4 address.
// Used by firewalls to allow single hosts access to a resource.
type Host struct {
	UID     int
	Name    string
	Address net.IP
	Comment string
}

// NewHost will return a pointer to a new host object.
// Will return an error if invalid IPv4 address in addr field.
func NewHost(name, addr, comment string) (*Host, error) {
	address := net.ParseIP(addr)
	if address == nil {
		return nil, fmt.Errorf("invalid host address: %s", addr)
	}
	return &Host{
		UID:     0,
		Name:    name,
		Address: address,
		Comment: comment,
	}, nil
}

// Match will return true if addr matches the host's address.
// Returns false if invalid IPv4 address - might be better to return an error?
func (h *Host) Match(addr *Host) bool {
	return h.Address.Equal(addr.Address)
}
