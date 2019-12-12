package core

import (
	"errors"
	"fmt"
	"net"
	"regexp"
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

// Match will return true if passed a network address that matches the object's address.
// The format of the address is [network address]/[subnet prefix] i.e. 192.168.1.0/24
// This will return false if incorrect format is used - maybe better to return error?
func (n *Network) Match(addr string) bool {
	// TODO: make this check more specific. Needs to match against [ipaddress]/[netmask]
	valid, err := regexp.MatchString(".*/.*", addr)
	if err != nil {
		return false
	}
	if !valid {
		return false
	}

	if n.Address.String() == addr {
		return true
	}
	return false
}

func (n *Network) containsHost(hostAddr string) (bool, error) {
	host := net.ParseIP(hostAddr)
	if host == nil {
		return false, fmt.Errorf("invalid host address: %s", hostAddr)
	}

	return n.Address.Contains(host), nil
}

func (n *Network) containsRange(rangeAddr string) (bool, error) {
	components, err := checkRangeFmt(rangeAddr)
	if err != nil {
		return false, fmt.Errorf("invalid range address: %v", err)
	}
	start := components[0]
	end := components[1]

	startIP := net.ParseIP(start)
	if startIP == nil {
		return false, fmt.Errorf("invalid host address: %s", start)
	}
	endIP := net.ParseIP(end)
	if endIP == nil {
		return false, fmt.Errorf("invalid host address: %s", end)
	}

	if valid := checkValidRange(startIP, endIP); !valid {
		return false, fmt.Errorf("start of range must be less than the end address: %s", rangeAddr)
	}

	return (n.Address.Contains(startIP) && n.Address.Contains(endIP)), nil
}

func (n *Network) containsNetwork(networkAddr string) (bool, error) {
	return false, errNotImplemented
}
