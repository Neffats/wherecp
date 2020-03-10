package core

import (
	"fmt"
	"strings"
)

const (
	TCP = iota
	UDP
	ICMP
	ARP
	IP
)

var proto = [5]string{"tcp", "udp", "icmp", "arp", "ip"}

// String2Proto returns the enum of string provided.
func String2Proto(protocol string) int {
	lowerProto := strings.ToLower(protocol)
	for i, p := range proto {
		if lowerProto == p {
			return i
		}
	}
	return -1
}

// Proto2String returns the string of the protocol enum.
func Proto2String(protocol int) string {
	if protocol < len(proto) || protocol > len(proto) {
		return ""
	}
	return proto[protocol]
}

type Port struct {
	UID      int
	Name     string
	Number   uint
	Protocol int
	Comment  string
}

func NewPort(name string, number uint, protocol, comment string) (*Port, error) {
	protoEnum := String2Proto(protocol)
	if protoEnum == -1 {
		return nil, fmt.Errorf("failed to create new Port object because invalid protocol provided: %s", protocol)
	}
	return &Port{
		Name:     name,
		Number:   number,
		Protocol: protoEnum,
		Comment:  comment,
	}, nil
}

func (p *Port) Value() (start uint, end uint, proto int) {
	start = p.Number
	end = p.Number
	proto = p.Protocol
	return
}

func (p *Port) Match(prt *Port) bool {
	return prt.Number == p.Number && prt.Protocol == prt.Protocol
}

func (p *Port) Contains(obj PortObject) bool {
	otherStart, otherEnd, otherProto := obj.Value()
	return otherStart == p.Number && otherEnd == p.Number && otherProto == p.Protocol
}
