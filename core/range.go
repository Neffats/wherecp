package core

import (
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"
)

// Range represents a range of IPv4 addresses.
// The start address must be smaller than the end address.
type Range struct {
	UID          int
	Name         string
	StartAddress net.IP
	EndAddress   net.IP
	Comment      string
}

// NewRange returns a pointer to a range object.
// start and end represent the start and end of the address range.
// Address format is the same as host i.e. 192.168.1.1
// Returns an error if start address is greater than the end address or an invalid address format.
func NewRange(name, start, end, comment string) (*Range, error) {
	r := new(Range)

	rangeStart := net.ParseIP(start)
	if rangeStart == nil {
		return nil, fmt.Errorf("invalid start address: %s", start)
	}
	rangeEnd := net.ParseIP(end)
	if rangeEnd == nil {
		return nil, fmt.Errorf("invalid start address: %s", start)
	}

	if valid := checkValidRange(rangeStart, rangeEnd); !valid {
		return r, fmt.Errorf("range start address must be less than the end address: %s-%s", start, end)
	}
	r.UID = 0
	r.Name = name
	r.StartAddress = rangeStart
	r.EndAddress = rangeEnd
	r.Comment = comment

	return r, nil
}

// Match will return true if the passed in range object's address matches.
func (r *Range) Match(addr *Range) bool {
	return reflect.DeepEqual(addr.StartAddress, r.StartAddress) && reflect.DeepEqual(addr.EndAddress, r.EndAddress)
}

// Checks if the format of the range string is valid.
func checkRangeFmt(addr string) ([]string, error) {
	// TODO: make this check more specific. Needs to match against [ipaddress]/[netmask]
	valid, err := regexp.MatchString(".*-.*", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to pattern match range address: %v", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid range address: %s", addr)
	}

	components := strings.Split(addr, "-")
	if len(components) != 2 {
		return nil, fmt.Errorf("range split failed wanted: %d, got: %d", 2, len(components))
	}
	return components, nil
}

func checkValidRange(start, end net.IP) bool {
	// The start of a range needs to be smaller than the end of it.
	s := ip2int(start)
	e := ip2int(end)
	if s > e {
		return false
	}
	return true
}

// Stole from https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func mask2int(mask net.IPMask) uint32 {
	if len(mask) == 16 {
		return binary.BigEndian.Uint32(mask[12:16])
	}
	return binary.BigEndian.Uint32(mask)
}

func (r *Range) containsHost(host *Host) (bool, error) {
	h := ip2int(host.Address)
	s := ip2int(r.StartAddress)
	e := ip2int(r.EndAddress)

	if h >= s && h <= e {
		return true, nil
	}
	return false, nil
}

// Only returns true if both the foreign range is the same or inside the self range.
func (r *Range) containsRange(foreignRange *Range) (bool, error) {
	s := ip2int(r.StartAddress)
	e := ip2int(r.EndAddress)

	fStart := ip2int(foreignRange.StartAddress)
	fEnd := ip2int(foreignRange.EndAddress)

	if fStart >= s && fStart <= e && fEnd >= s && fEnd <= e {
		return true, nil
	}

	return false, nil
}

func (r *Range) containsNetwork(network *Network) (bool, error) {
	return false, errNotImplemented
}
