package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/Neffats/ip"
)

// Range represents a range of IPv4 addresses.
// The start address must be smaller than the end address.
type Range struct {
	uid          string
	name         string
	startAddress *ip.Address
	endAddress   *ip.Address
	comment      string
}

// NewRange returns a pointer to a range object.
// start and end represent the start and end of the address range.
// Address format is the same as host i.e. 192.168.1.1
// Returns an error if start address is greater than the end address or an invalid address format.
func NewRange(name, start, end, comment string) (*Range, error) {
	r := new(Range)

	rangeStart, err := ip.NewAddress(start)
	if err != nil {
		return nil, fmt.Errorf("invalid start address: %v", err)
	}
	rangeEnd, err := ip.NewAddress(end)
	if err != nil {
		return nil, fmt.Errorf("invalid end address: %v", err)
	}

	if *rangeStart > *rangeEnd {
		return r, fmt.Errorf("range start address must be less than the end address: %s-%s", start, end)
	}
	uid := uuid.New()
	r.uid = uid.String()
	r.name = name
	r.startAddress = rangeStart
	r.endAddress = rangeEnd
	r.comment = comment

	return r, nil
}

func (r *Range) Unpack() []NetworkObject {
	result := make([]NetworkObject, 0)
	result = append(result,
		NetworkObject{
			Start: *r.startAddress,
			End: *r.endAddress,
		})
	return result
}

// Match will return true if the passed in range object's address matches.
func (r *Range) Match(addr *Range) bool {
	return *r.startAddress == *addr.startAddress && *r.endAddress == *addr.endAddress
}

// Contains will return true if obj is contained by the range.
func (r *Range) Contains(obj NetworkUnpacker) bool {
	compare := obj.Unpack()
	for _, c := range compare {
		if c.Start < *r.startAddress || c.End > *r.endAddress {
			return false
		}
	}
	return true
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
