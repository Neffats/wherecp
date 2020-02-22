package core

import (
	"fmt"

	"github.com/Neffats/ip"
)

// Host represents a single IPv4 host object, a single IPv4 address.
// Used by firewalls to allow single hosts access to a resource.
type Host struct {
	UID     int
	Name    string
	Address *ip.Address
	Comment string
}

// NewHost will return a pointer to a new host object.
// Will return an error if invalid IPv4 address in addr field.
func NewHost(name, addr, comment string) (*Host, error) {
	address, err := ip.NewAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid host address: %v", err)
	}
	return &Host{
		UID:     0,
		Name:    name,
		Address: address,
		Comment: comment,
	}, nil
}

func (h *Host) Value() (start *ip.Address, end *ip.Address) {
	start = h.Address
	end = h.Address
	return
}

// Match will return true if addr matches the host's address.
// Returns false if invalid IPv4 address - might be better to return an error?
func (h *Host) Match(obj NetworkObject) bool {
	start, end := obj.Value()
	if *start == *h.Address && *end == *h.Address {
		return true
	}
	return false
}
