package core

import (
	"errors"
	"fmt"

	"github.com/Neffats/ip"
	"github.com/google/uuid"
)

var (
	errNotImplemented = errors.New("not implemented")
)

const (
	// Max value of an IPv4 address.
	addrMax = ip.Address(4294967295)
)

// Network represents an IPv4 subnet.
// Used by firewalls to allow whole networks access to a resource.
type Network struct {
	uid     string
	name    string
	address *ip.Address
	mask    *ip.Address
	comment string
}

// NewNetwork returns a ptr to a new Network object.
// addr is the network ip address for the subnet i.e. 192.168.1.0.
// mask is the subnet mask for the network i.e. 255.255.255.0
func NewNetwork(name, addr, mask, comment string) (*Network, error) {
	network := new(Network)

	netAddr, err := ip.NewAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create network address: %v", err)
	}
	netMask, err := ip.NewAddress(mask)
	if err != nil {
		return nil, fmt.Errorf("failed to create mask address: %v", err)
	}

	// Check if addr is actually the network address for the subnet.
	// Address shouldn't change if (bitwise) and'd with netmask.
	if *ip.Mask(netAddr, netMask) != *netAddr {
		return nil, fmt.Errorf("address not the network address for supplied subnet: %s/%s", addr, mask)
	}
	uid := uuid.New()
	network.uid = uid.String()
	network.name = name
	network.address = netAddr
	network.mask = netMask
	network.comment = comment

	return network, nil
}

// Value returns the first and last address in the Network's Address range (network and broadcast).
// Statisfies the NetworkObject interface.
func (n *Network) Unpack() []NetworkObject {
	// get inverse of the subnet mask
	invMask := *n.mask ^ addrMax

	start := *n.address
	// Or the network address with the inverse of the mask to get the last address in the subnet.
	endAddr := *n.address | invMask
	end := endAddr
	result := make([]NetworkObject, 0)
	result = append(result,
		NetworkObject{
			Start: start,
			End: end,
		})

	return result
}

// Match will return true if passed a network that has a matching address.
func (n *Network) Match(addr *Network) bool {
	return *n.address == *addr.address && *n.mask == *addr.mask
}

// Contains takes a NetworkObject, returns true if the object's start and end Address
// falls within the network Address range.
func (n *Network) Contains(obj NetworkUnpacker) bool {
	compare := obj.Unpack()
	this := n.Unpack()
	self := this[0]
	for _, c := range compare {
		if c.Start < self.Start || c.End > self.End {
			return false
		}
	}
	return true
}
