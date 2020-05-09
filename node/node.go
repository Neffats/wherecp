package node

import (
	"github.com/Neffats/ip"
)

type Subnet struct {
	Address *ip.Address
	Mask    *ip.Address
}

type Node struct {
	ConnectedNetworks []Subnet
	Rules RuleStorer
	Hosts HostStorer
	Networks NetworkStorer
	Ranges RangeStorer
	Groups GroupStorer
}
