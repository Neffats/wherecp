package core

import (
	"bytes"
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

	if valid := checkValidRange(rangeStart, rangeStart); !valid {
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
	comp := bytes.Compare(start, end)
	if comp != -1 {
		return false
	}
	return true
}

func convertIP(addr net.IP) int {
	final := 0
	for _, b := range addr {
		final <<= 8
		final &= int(b)
	}
	return final
}

func (r *Range) containsHost(host *Host) (bool, error) {
	return false, errNotImplemented
}

func (r *Range) containsRange(rng *Range) (bool, error) {
	return false, errNotImplemented
}

func (r *Range) containsNetwork(network *Network) (bool, error) {
	return false, errNotImplemented
}
