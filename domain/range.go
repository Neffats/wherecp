package core

import (
	"bytes"
	"errors"
	"net"
	"reflect"
	"regexp"
	"strings"
)

type Range struct {
	Uid          int
	Name         string
	StartAddress net.IP
	EndAddress   net.IP
	Comment      string
}

func NewRange(name, start, end, comment string) (*Range, error) {
	r := new(Range)

	rangeStart := net.ParseIP(start)
	rangeEnd := net.ParseIP(end)

	// The start of a range needs to be smaller than the end of it.
	comp := bytes.Compare(rangeStart, rangeEnd)
	if comp != -1 {
		return r, errors.New("range start address must be less than the end address")
	}
	r.Uid = 0
	r.Name = name
	r.StartAddress = rangeStart
	r.EndAddress = rangeEnd
	r.Comment = comment

	return r, nil
}

// Match will return true if the range string matches the range objects.
// Will return false if format of addr is incorrect.
// Correct format == [ip address]-[ip address] i.e 192.168.1.1-192.168.1.10.
// The first address in the range must be smaller than the second.
func (r *Range) Match(addr string) bool {
	// TODO: make this check more specific. Needs to match against [ipaddress]/[netmask]
	valid, err := regexp.MatchString(".*-.*", addr)
	if err != nil {
		return false
	}
	if !valid {
		return false
	}

	components := strings.Split(addr, "-")
	if len(components) != 2 {
		return false
	}
	start := components[0]
	end := components[1]

	startip := net.ParseIP(start)
	endip := net.ParseIP(end)

	if reflect.DeepEqual(startip, r.StartAddress) && reflect.DeepEqual(endip, r.EndAddress) {
		return true
	}
	return false
}
