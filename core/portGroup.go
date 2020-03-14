package core

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// PortGroup groups together different Port, PortRanges and other PortGroup objects.
type PortGroup struct {
	UID     int
	Name    string
	Ports   []*Port
	Ranges  []*PortRange
	Groups  []*PortGroup
	Comment string
}

// NewPortGroup returns a new empty oort group.
func NewPortGroup(name, comment string) *PortGroup {
	return &PortGroup{
		UID:     0,
		Name:    name,
		Ports:   make([]*Port, 0),
		Ranges:  make([]*PortRange, 0),
		Groups:  make([]*PortGroup, 0),
		Comment: comment,
	}
}

// Add will add the specified object to the group.
// Supported types: Port/Port Range/Port Group
func (pg *PortGroup) Add(obj interface{}) error {
	present, err := pg.HasObject(obj)
	if err != nil {
		return fmt.Errorf("failed to check if object is already a group member: %v", err)
	}
	if present {
		return fmt.Errorf("object is already a member of this group: %s", obj)
	}

	switch v := obj.(type) {
	case *Port:
		pg.addPort(v)
	case *PortRange:
		pg.addPortRange(v)
	case *PortGroup:
		pg.addPortGroup(v)
	default:
		return errors.New("unsupported data type")
	}
	return nil
}

func (pg *PortGroup) addPort(p *Port) {
	// Ordered smallest to largets, first by port number then by protocol number.
	i := sort.Search(len(pg.Ports), func(i int) bool {
		return pg.Ports[i].Number > p.Number && pg.Ports[i].Protocol >= p.Protocol
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newPorts := make([]*Port, len(pg.Ports)+1)
	// Shift the slice forward by one at the insert location.
	copy(newPorts[:i], pg.Ports[:i])
	copy(newPorts[i+1:], pg.Ports[i:])
	// Append port at the insert location.
	newPorts[i] = p
	pg.Ports = newPorts
}

func (pg *PortGroup) addPortRange(pr *PortRange) {
	// Ordered smallest to largest by start range start (first port) first
	// then by range end (last port) second, then by protocol number. Smaller ranges will come before
	// larger ranges i.e. 1-2, 1-5, 2-4
	i := sort.Search(len(pg.Ranges), func(i int) bool {
		thisStart, thisEnd, thisProto := pg.Ranges[i].Value()
		otherStart, otherEnd, otherProto := pr.Value()

		start := thisStart >= otherStart
		end := thisEnd >= otherEnd
		proto := thisProto >= otherProto
		return start && end && proto
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newRanges := make([]*PortRange, len(pg.Ranges)+1)
	// Shift the slice forward by one at the insert location.
	copy(newRanges[:i], pg.Ranges[:i])
	copy(newRanges[i+1:], pg.Ranges[i:])
	// Append port at the insert location.
	newRanges[i] = pr
	pg.Ranges = newRanges
}

func (pg *PortGroup) addPortGroup(grp *PortGroup) {
	// Ordered alphabetically by Group name.
	i := sort.Search(len(pg.Groups), func(i int) bool {
		return pg.Groups[i].Name >= grp.Name
	})

	// TODO: Is there a nicer way of doing this?
	// Create a new bigger slice.
	newGroup := make([]*PortGroup, len(pg.Groups)+1)
	// Shift the slice forward by one at the insert location.
	copy(newGroup[:i], pg.Groups[:i])
	copy(newGroup[i+1:], pg.Groups[i:])
	// Append group at the insert location.
	newGroup[i] = grp
	pg.Groups = newGroup
}

// HasObject returns true if the group has a members object whose type and address matches the supplied object.
func (pg *PortGroup) HasObject(obj interface{}) (bool, error) {
	// TODO: Make more efficient since lists are now ordered.
	switch v := obj.(type) {
	case *Port:
		if len(pg.Ports) < 1 {
			return false, nil
		}
		var i int
		// Edge case handling. When len() == 0, sort.Search() was returning an index of 1 which is oob.
		if len(pg.Ports) == 1 {
			i = 0
		} else {
			i = sort.Search(len(pg.Ports), func(i int) bool {
				keySt, keyEnd, keyProto := v.Value()
				midSt, midEnd, midProto := pg.Ports[i].Value()

				return keySt == midSt && keyEnd == midEnd && keyProto == midProto
			})
		}

		// Check that what we go makes sense.
		if i == -1 || i >= len(pg.Ports) {
			return false, nil
		}

		// Double check that objects match.
		if pg.Ports[i].Match(v) {
			return true, nil
		}
	case *PortRange:
		var i int
		// Edge case handling. When len() == 0, sort.Search() was return index of 1 with is oob.
		if len(pg.Ranges) == 1 {
			i = 0
		} else {
			i = sort.Search(len(pg.Ranges), func(i int) bool {
				keySt, keyEnd, keyProto := v.Value()
				midSt, midEnd, midProto := pg.Ranges[i].Value()

				return keySt == midSt && keyEnd == midEnd && keyProto == midProto
			})
		}

		// Check that what we go makes sense.
		if i == -1 || i >= len(pg.Ranges) {
			return false, nil
		}

		// Double check that objects match.
		if pg.Ranges[i].Match(v) {
			return true, nil
		}
	case *PortGroup:
		var i int
		// Edge case handling. When len() == 0, sort.Search() was return index of 1 with is oob.
		if len(pg.Groups) == 1 {
			i = 0
		} else {
			i = sort.Search(len(pg.Groups), func(i int) bool {
				return pg.Groups[i].Name == v.Name
			})
		}

		// Check that what we go makes sense.
		if i == -1 || i >= len(pg.Groups) {
			return false, nil
		}

		// Double check that objects match.
		if pg.Groups[i].Match(v) {
			return true, nil
		}
	default:
		return false, errors.New("unsupported data type")
	}

	// Check if any of it's group members contain the object.
	for _, grp := range pg.Groups {
		has, err := grp.HasObject(obj)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

// Match will return true if both groups are identical.
func (pg *PortGroup) Match(grp *PortGroup) bool {
	return reflect.DeepEqual(pg, grp)
}

// MatchContent returns true if both groups contain the exact same members.
func (pg *PortGroup) MatchContent(grp *PortGroup) bool {
	// Check if the lengths of the groups match.
	// If they don't then the two groups must be different.
	if len(pg.Ports) != len(grp.Ports) {
		return false
	}
	if len(pg.Ranges) != len(grp.Ranges) {
		return false
	}
	if len(pg.Groups) != len(grp.Groups) {
		return false
	}

	var match bool

	// Compare Hosts of groups.
	// All group members are sorted, so all members should be in the same location.
	for i := 0; i < len(pg.Ports); i++ {
		match = pg.Ports[i].Match(grp.Ports[i])
		if !match {
			return false
		}
	}

	// Compare Ranges of both groups.
	for i := 0; i < len(pg.Ranges); i++ {
		match = pg.Ranges[i].Match(grp.Ranges[i])
		if !match {
			return false
		}
	}

	// Compare Groups of groups.
	for i := 0; i < len(pg.Groups); i++ {
		match = pg.Groups[i].Match(grp.Groups[i])
		if !match {
			return false
		}
	}

	return true
}
