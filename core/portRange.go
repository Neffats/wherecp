package core

import "fmt"

type PortRange struct {
	uid      int
	name     string
	start    uint
	end      uint
	protocol int
	comment  string
}

func NewPortRange(name string, start, end uint, protocol, comment string) (*PortRange, error) {
	protoEnum := String2Proto(protocol)
	if protoEnum == -1 {
		return nil, fmt.Errorf("failed to create new Port object because invalid protocol provided: %s", protocol)
	}
	return &PortRange{
		name:     name,
		start:    start,
		end:      end,
		protocol: protoEnum,
		comment:  comment,
	}, nil
}

func (pr *PortRange) Value() (start uint, end uint, proto int) {
	start = pr.start
	end = pr.end
	proto = pr.protocol
	return
}

func (pr *PortRange) Match(other *PortRange) bool {
	return other.start == pr.start && other.end == pr.end && other.protocol == pr.protocol
}

func (pr *PortRange) Contains(obj PortObject) bool {
	otherStart, otherEnd, otherProto := obj.Value()
	return pr.start <= otherStart && pr.end >= otherEnd && otherProto == pr.protocol
}
