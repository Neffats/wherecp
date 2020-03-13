package core

import "fmt"

type PortRange struct {
	UID      int
	Name     string
	Start    uint
	End      uint
	Protocol int
	Comment  string
}

func NewPortRange(name string, start, end uint, protocol, comment string) (*PortRange, error) {
	protoEnum := String2Proto(protocol)
	if protoEnum == -1 {
		return nil, fmt.Errorf("failed to create new Port object because invalid protocol provided: %s", protocol)
	}
	return &PortRange{
		Name:     name,
		Start:    start,
		End:      end,
		Protocol: protoEnum,
		Comment:  comment,
	}, nil
}

func (pr *PortRange) Value() (start uint, end uint, proto int) {
	start = pr.Start
	end = pr.End
	proto = pr.Protocol
	return
}

func (pr *PortRange) Match(other *PortRange) bool {
	return other.Start == pr.Start && other.End == pr.End && other.Protocol == pr.Protocol
}

func (pr *PortRange) Contains(obj PortObject) bool {
	otherStart, otherEnd, otherProto := obj.Value()
	return pr.Start <= otherStart && pr.End >= otherEnd && otherProto == pr.Protocol
}
