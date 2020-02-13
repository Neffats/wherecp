package core

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"

	"github.com/Neffats/ip"
)

var (
	errNotImplemented = errors.New("not implemented")
)

const (
	addrMax = ip.Address(4294967295)
)

// Network represents an IPv4 subnet.
// Used by firewalls to allow whole networks access to a resource.
type Network struct {
	UID     int
	Name    string
	Address *ip.Address
	Mask    *ip.Address
	Comment string
}

// NewNetwork returns a ptr to a new Network object.
// addr is the network ip address for the subnet i.e. 192.168.1.0.
// mask is the network prefix for the netmask i.e. 24 = 255.255.255.0
// Will return an error if the netmask is not valid i.e. > 32 or < 0 or invalid network address.
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
	// Address shouldn't change if (bitwise) anded with netmask.
	if ip.Mask(netAddr, netMask) != netAddr {
		return nil, fmt.Errorf("address not the network address for supplied subnet: %s/%s", addr, mask)
	}
	network.UID = 0
	network.Name = name
	network.Address = netAddr
	network.Mask = netMask
	network.Comment = comment

	return network, nil
}

func (n *Network) Value() (start *ip.Address, end *ip.Address) {
	// get inverse of the subnet mask
	invMask := n.Mask ^ addrMax

	start := n.Address
	// Or the network address with the inverse of the mask to get the last address in the subnet.
	end := n.Address | invMask
	return
}

// Match will return true if passed a network that has a matching address.
func (n *Network) Match(addr *Network) bool {
	return reflect.DeepEqual(n.Address, addr.Address)
}

func (n *Network) Contains(obj NetworkObject) bool {
	compStart, compEnd := obj.Value()
	thisStart, thisEnd := n.Value()
	if compStart > thisStart && compEnd < thisEnd {
		return true
	}
	return false
}

/*
// Contains will return true if the network contains the provided object.
// Only returns true if all of the object is inside the network.
// i.e. if a range starts inside and finishes outside it will return false.
func (n *Network) Contains(obj interface{}) (bool, error) {
	var (
		result bool
		err    error
	)
	switch v := obj.(type) {
	case *Host:
		result, err = n.containsHost(v)
	case *Network:
		result, err = n.containsNetwork(v)
	case *Range:
		result, err = n.containsRange(v)
	default:
		return false, errors.New("provided data type is unsupported")
	}

	return result, err
}
*/

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
