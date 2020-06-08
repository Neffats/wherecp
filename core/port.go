package core

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
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
	if protocol < 0 || protocol > len(proto)-1 {
		return ""
	}
	return proto[protocol]
}

type Port struct {
	uid      string
	name     string
	number   uint
	protocol int
	comment  string
}

func NewPort(name string, number uint, protocol, comment string) (*Port, error) {
	protoEnum := String2Proto(protocol)
	if protoEnum == -1 {
		return nil, fmt.Errorf("failed to create new Port object because invalid protocol provided: %s", protocol)
	}
	uid := uuid.New()
	return &Port{
		uid: uid.String(),
		name:     name,
		number:   number,
		protocol: protoEnum,
		comment:  comment,
	}, nil
}

func (p *Port) Value() (start uint, end uint, proto int) {
	start = p.number
	end = p.number
	proto = p.protocol
	return
}

func (p *Port) Match(prt *Port) bool {
	return prt.number == p.number && prt.protocol == p.protocol
}

func (p *Port) Contains(obj PortObject) bool {
	otherStart, otherEnd, otherProto := obj.Value()
	return otherStart == p.number && otherEnd == p.number && otherProto == p.protocol
}
