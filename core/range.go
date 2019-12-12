package core

import (
	"bytes"
	"errors"
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
		return r, errors.New("range start address must be less than the end address")
	}
	r.UID = 0
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
	components, err := checkRangeFmt(addr)
	if err != nil {
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

func (r *Range) containsHost(hostAddr string) (bool, error) {
	return false, errNotImplemented
}

func (r *Range) containsRange(rangeAddr string) (bool, error) {
	return false, errNotImplemented
}

func (r *Range) containsNetwork(networkAddr string) (bool, error) {
	return false, errNotImplemented
}
