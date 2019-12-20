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
	netMask := net.CIDRMask(n, 32)

	// Check if addr is actually the network address for the subnet.
	// Address shouldn't change if (bitwise) anded with netmask.
	if (ip2int(netAddr) & mask2int(netMask)) != ip2int(netAddr) {
		return network, fmt.Errorf("address not the network address for supplied subnet: %s/%s", addr, mask)
	}
	network.UID = 0
	network.Name = name
	network.Address = net.IPNet{
		IP:   netAddr,
		Mask: netMask,
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

func (n *Network) containsNetwork(foreignN *Network) (bool, error) {
	if ip2int(n.Address.IP) <= ip2int(foreignN.Address.IP) {
		if ip2int(n.broadcast()) >= ip2int(foreignN.broadcast()) {
			return true, nil
		}
	}
	return false, nil
}

func (n *Network) broadcast() net.IP {
	broadcast := net.IP(make([]byte, 4))
	for i := range n.Address.IP[12:16] {
		broadcast[i] = n.Address.IP[12+i] | ^n.Address.Mask[i]
	}

	return broadcast
}
