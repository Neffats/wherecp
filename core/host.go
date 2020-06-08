package core

import (
	"fmt"

	"github.com/Neffats/ip"
	"github.com/google/uuid"
)

// Host represents a single IPv4 host object, a single IPv4 address.
// Used by firewalls to allow single hosts access to a resource.
type Host struct {
	uid     string
	name    string
	address *ip.Address
	comment string
}

// NewHost will return a pointer to a new host object.
// Will return an error if invalid IPv4 address in addr field.
func NewHost(name, addr, comment string) (*Host, error) {
	address, err := ip.NewAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid host address: %v", err)
	}
	uid := uuid.New()
	return &Host{
		uid:     uid.String(),
		name:    name,
		address: address,
		comment: comment,
	}, nil
}

func (h *Host) UID() string {
	return h.uid
}

func (h *Host) Unpack() []NetworkObject {
	result := make([]NetworkObject, 0)
	result = append(result,
		NetworkObject{
			Start: *h.address,
			End: *h.address,
		})
	return result
}

func (h *Host) Match(obj *Host) bool {
	return *h.address == *obj.address
}

// Contains will return true if addr matches the host's address.
// Returns false if invalid IPv4 address - might be better to return an error?
func (h *Host) Contains(obj NetworkUnpacker) bool {
	networks := obj.Unpack()
	for _, n := range networks {
		if n.Start != *h.address || n.End != *h.address {
			return false
		}
	}
	return true
}
