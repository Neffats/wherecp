package core

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
)

type Network struct {
	Uid     int
	Name    string
	Address net.IPNet
	Netmask net.IPMask
	Comment string
}

// NewNetwork returns a ptr to a new Network object.
// addr is the network ip address for the subnet i.e. 192.168.1.0.
// mask is the network prefix for the netmask i.e. 24 = 255.255.255.0
// Will return an error if the netmask is not valid i.e. > 32 or < 0
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
	network.Uid = 0
	network.Name = name
	network.Address = net.IPNet{
		IP:   net.ParseIP(addr),
		Mask: net.CIDRMask(n, 32),
	}
	network.Comment = comment

	return network, nil
}

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
