package core

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
)

var (
	errNotImplemented = errors.New("not implemented")
)

// Network represents an IPv4 subnet.
// Used by firewalls to allow whole networks access to a resource.
type Network struct {
	UID     int
	Name    string
	Address net.IPNet
	Netmask net.IPMask
	Comment string
}

// NewNetwork returns a ptr to a new Network object.
// addr is the network ip address for the subnet i.e. 192.168.1.0.
// mask is the network prefix for the netmask i.e. 24 = 255.255.255.0
// Will return an error if the netmask is not valid i.e. > 32 or < 0 or invalid network address.
func NewNetwork(name, addr, mask, comment string) (*Network, error) {
	network := new(Network)
	n, err := strconv.Atoi(mask)
	if err != nil {
		return network, fmt.Errorf("failed to convert mask to int: %v", err)
	}
	// Make sure that we got a valid network prefix.
	if n < 0 || n > 32 {
		return network, fmt.Errorf("invalid mask provided: %v", n)
	}

	netAddr := net.ParseIP(addr)
	if netAddr == nil {
		return network, fmt.Errorf("invalid network address: %s", addr)
	}
	network.UID = 0
	network.Name = name
	network.Address = net.IPNet{
		IP:   netAddr,
		Mask: net.CIDRMask(n, 32),
	}
	network.Comment = comment

	return network, nil
}

// Match will return true if passed a network that has a matching address.
func (n *Network) Match(addr *Network) bool {
	return reflect.DeepEqual(n.Address, addr.Address)
}

func (n *Network) containsHost(h *Host) (bool, error) {
	return n.Address.Contains(h.Address), nil
}

func (n *Network) containsRange(r *Range) (bool, error) {
	return (n.Address.Contains(r.StartAddress) && n.Address.Contains(r.EndAddress)), nil
}

func (n *Network) containsNetwork(network *Network) (bool, error) {
	return false, errNotImplemented
}
